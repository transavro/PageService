package apihandler

import (
	pb "PageService/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func insertPageDB(ctx context.Context, pageColl *mongo.Collection, page *pb.Page) ( err error) {
	page.PageId = pageIdGenerator()
	page.CreatedAt = ptypes.TimestampNow()
	_, err  = pageColl.InsertOne(ctx, page)
	return err
}

func getPageDB(ctx context.Context, pageColl *mongo.Collection, pageId string) (*pb.Page, error)  {
	var (
		page *pb.Page
		err error
	)
	if err = pageColl.FindOne(ctx, bson.M{"pageid" : pageId}).Decode(&page); err != nil {
		return nil, err
	}else {
		return page, err
	}
}

func updatePageDB(ctx context.Context, pageColl *mongo.Collection, pageId string, page *pb.Page) error {
	page.UpdatedAt = ptypes.TimestampNow()
	return pageColl.FindOneAndReplace(ctx, bson.M{"pageid" : pageId}, page).Err()
}


func deletePageDB(ctx context.Context, pageColl *mongo.Collection, pageId string) (err error){
	return pageColl.FindOneAndDelete(ctx, bson.M{"pageid" : pageId}).Err()
}

func getAllDB(ctx context.Context, pageColl *mongo.Collection) ( cur *mongo.Cursor, err error)  {
	if cur, err = pageColl.Find(ctx, bson.M{}); err != nil {
		return nil, err
	}else {
		return cur, nil
	}
}



func distinctValueDB(ctx context.Context, tileColl *mongo.Collection, field string)(*pb.DropDownResp, error)  {
	if result, err := tileColl.Distinct(ctx, field, bson.D{}); err != nil {
			return nil,  status.Errorf(codes.Internal, "Filed %s Does not exists. [Fields are case Sensitive]", field)
	}else if len(result) == 0 {
		return nil,  status.Errorf(codes.Internal, "Filed %s Does not exists. [Fields are case Sensitive]", field)
	}else {
		temp := make([]string, len(result))
		for i, v := range result {
			temp[i] = fmt.Sprint(v)
		}
		return &pb.DropDownResp{Result: temp} , nil
	}
}