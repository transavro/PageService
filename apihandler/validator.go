package apihandler

import (
	pb "PageService/proto"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)


func pageValidator(ctx context.Context, tileColl *mongo.Collection, page *pb.Page) (err error) {
	if err = rowValidator(ctx, tileColl, page.Row); err != nil {
		return  err
	}else if err = carouselValidator(page.Carousel); err != nil {
		return  err
	}else {
		return nil
	}
}

func rowValidator(ctx context.Context, tileColl *mongo.Collection, rows []*pb.Row) error  {

	var checker []int32

	//check if the rowfilters are valid
	for  _ , row := range rows {

		// check rowType
		if row.RowType == pb.RowType_Editorial{
			if len(row.GetRowTileIds()) == 0 {
				return status.Errorf(codes.InvalidArgument, "Row = %s of type %s have empty rowTileArray", row.RowName, row.RowType.String())
			}
			continue
		}else if row.RowType == pb.RowType_Web{
			if v, ok := row.RowFilters["source"]; !ok && len(v.Values) == 0 {
				return status.Errorf(codes.InvalidArgument, "Row = %s of type %s must have 'source' as a rowFilter, with sourceName as the value.", row.RowName, row.RowType.String())
			}else  {
				for _, value := range v.Values {
					if "youtube" == strings.ToLower(value){
						if _, ok = row.RowFilters["playlist"]; !ok{
							if _, ok = row.RowFilters["channel"]; !ok {
								if _, ok = row.RowFilters["search"]; !ok {

									return status.Errorf(codes.InvalidArgument,
										"Row = %s of type %s, have 'source=youtube', " +
										"but not have 'playlist=PLAYLIST_ID' or 'channel=CHANNEL_ID' or 'search=QUERY' as rowFilter.",
										row.RowName, row.RowType.String())
								}
							}
						}
						break
					}
				}
			}
			continue
		}

		for k, v := range row.RowFilters {
			if !isMongoKeyExists(tileColl , k, ctx) {
				return status.Errorf(codes.InvalidArgument, "%s rowFilter key %s is invalid", row.RowName, k)
			}
			//trimming start and end spaces from values
			row.RowFilters[k] = &pb.RowFilterValue{Values: trimSpaces(v.Values)}
		}
		


		// check rowSort
		for k, v := range row.GetRowSort() {
			if v == 0 {
				return status.Errorf(codes.InvalidArgument, "%s RowSort %s=0 is invalid. -1 for sort descending and 1 for sort ascending.", row.RowName, k)
			}
			if !isMongoKeyExists(tileColl , k, ctx) {
				return status.Errorf(codes.InvalidArgument, "%s RowSort key %s is invalid", row.RowName, k)
			}
		}
		
		checker = append(checker, row.RowIndex)
	}

	// check row Index validation
	return indexChecker(checker)
}

func carouselValidator(carousels []*pb.Carousel) error {

	var checker []int32

	for _ , carousel := range carousels {
		checker = append(checker, carousel.Index)
	}

	// check Carousel Index validation
	return indexChecker(checker)
}
