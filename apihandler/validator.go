package apihandler

import(
	pb "PageService/proto"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		for k, v := range row.RowFilters {
			if !isMongoKeyExists(tileColl , k, ctx) {
				return status.Errorf(codes.InvalidArgument, "%s rowFilter key %s is invalid", row.RowName, k)
			}
			//trimming start and end spaces from values
			row.RowFilters[k] = &pb.RowFilterValue{Values: trimSpaces(v.Values)}
		}
		
		// check rowType
		if row.RowType == pb.RowType_Editorial{
			if len(row.GetRowTileIds()) == 0 {
				return status.Errorf(codes.InvalidArgument, "Row = %s of type %s have empty rowTileArray", row.RowName, row.RowType.String())
			}
		}

		// checl rowSort
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
