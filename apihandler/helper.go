package apihandler

import (
	pb "PageService/proto"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func pageResolver(ctx context.Context, tileColl *mongo.Collection, page *pb.Page) (*pb.ResultPage, error) {
	var (
		resultPage pb.ResultPage
		resultRow  *pb.ResultRow
		err        error
	)
	resultPage.Carousels = page.Carousel
	resultPage.PageId = page.PageId
	resultPage.PageName = page.PageName

	for _, row := range page.Row {
		if resultRow, err = rowResolver(ctx, tileColl, row); err != nil {
			return nil, err
		} else {
			resultPage.Rows = append(resultPage.Rows, resultRow)
		}
	}
	return &resultPage, err
}

func rowResolver(ctx context.Context, tileColl *mongo.Collection, row *pb.Row) (*pb.ResultRow, error) {
	switch row.RowType {
	case pb.RowType_Editorial:
		{
			return makeEditorialRow(ctx, tileColl, row)
		}
	case pb.RowType_Dynamic:
		{
			return makeDynamicRow(ctx, tileColl, row)
		}
	default:
		return nil, nil
	}
}

func makeEditorialRow(ctx context.Context, tileColl *mongo.Collection, row *pb.Row) (*pb.ResultRow, error) {
	if cur, err := tileColl.Aggregate(ctx, makeEditorialPL(row.RowTileIds)); err != nil {
		return nil, err
	} else {
		return makeDeliveryRow(ctx, cur, row)
	}
}

func makeEditorialPL(rowIds []string) mongo.Pipeline {
	// creating pipes for mongo aggregation for Editorial Row
	stages := mongo.Pipeline{}
	stages = append(stages, bson.D{{"$match", bson.M{"refid": bson.M{"$in": rowIds}}}})
	stages = append(stages, bson.D{{"$lookup", bson.M{"from": "optimus_monetize", "localField": "refid", "foreignField": "refid", "as": "play"}}})
	stages = append(stages, bson.D{{"$replaceRoot", bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{bson.M{"$arrayElemAt": bson.A{"$play", 0}}, "$$ROOT"}}}}}) //adding stage 3  ==> https://docs.mongodb.com/manual/reference/operator/aggregation/mergeObjects/#exp._S_mergeObjects
	stages = append(stages, bson.D{{"$project", bson.M{"play": 0}}})
	stages = append(stages, bson.D{{"$project", bson.M{"_id": 0,
		"title":        "$metadata.title",
		"poster":       "$media.landscape",
		"portriat":     "$media.portrait",
		"video":        "$media.video",
		"type":         "$tiletype",
		"isDetailPage": "$content.detailpage",
		"contentId":    "$refid",
		"play":         "$contentavailable",
	}}})
	return stages
}

func makeDynamicRow(ctx context.Context, tileColl *mongo.Collection, row *pb.Row) (*pb.ResultRow, error) {
	if cur, err := tileColl.Aggregate(ctx, makeDynamicPL(row)); err != nil {
		return nil, err
	} else {
		return makeDeliveryRow(ctx, cur, row)
	}
}

func makeDynamicPL(row *pb.Row) mongo.Pipeline {

	stages := mongo.Pipeline{}
	var queryArray []bson.E

	//Row Filter
	for key, value := range row.GetRowFilters() {
		if value.GetValues() != nil && len(value.GetValues()) > 0 {
			queryArray = append(queryArray, bson.E{Key: key, Value: bson.M{"$in": value.GetValues()}})
		}
	}
	stages = append(stages, bson.D{{"$match", queryArray}})

	//Row Sort
	//if len(row.GetRowSort()) > 0 {
	//	queryArray = queryArray[:0]
	//	for key, value := range row.GetRowSort() {
	//		queryArray = append(queryArray, bson.E{Key: key, Value: value})
	//	}
	//	stages = append(stages, bson.D{{"$sort", queryArray}})
	//}

	// to make our delivery schema
	stages = append(stages, bson.D{{"$limit", 150}})
	stages = append(stages, bson.D{{"$lookup", bson.M{"from": "optimus_monetize", "localField": "refid", "foreignField": "refid", "as": "play"}}})
	stages = append(stages, bson.D{{"$replaceRoot", bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{bson.M{"$arrayElemAt": bson.A{"$play", 0}}, "$$ROOT"}}}}}) //adding stage 3  ==> https://docs.mongodb.com/manual/reference/operator/aggregation/mergeObjects/#exp._S_mergeObjects
	stages = append(stages, bson.D{{"$project", bson.M{"play": 0}}})
	stages = append(stages, bson.D{{"$project", bson.M{"_id": 0,
		"title":        "$metadata.title",
		"poster":       "$media.landscape",
		"portriat":     "$media.portrait",
		"video":        "$media.video",
		"type":         "$tiletype",
		"isDetailPage": "$content.detailpage",
		"contentId":    "$refid",
		"play":         "$contentavailable",
	}}})

	return stages
}

func makeDeliveryRow(ctx context.Context, cur *mongo.Cursor, row *pb.Row) (*pb.ResultRow, error) {
	var (
		resultRow pb.ResultRow
		err       error
	)
	resultRow.RowName = row.RowName
	resultRow.RowIndex = row.RowIndex

	for cur.Next(ctx) {
		var content *pb.Content
		if err = cur.Decode(&content); err != nil {
			return nil, err
		} else {
			resultRow.Tiles = append(resultRow.Tiles, content)
		}
	}
	if err = cur.Close(ctx); err != nil {
		return nil, err
	} else {
		return &resultRow, err
	}
}
