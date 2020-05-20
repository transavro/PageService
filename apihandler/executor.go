package apihandler

import (
	pb "PageService/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Server struct {
	PageCollection  *mongo.Collection
	TilesCollection *mongo.Collection
}

func (s Server) CreatePage(ctx context.Context, page *pb.Page) (*pb.Page, error) {
	var err error
	if err = pageValidator(ctx, s.TilesCollection, page); err != nil {
		return nil, err
	} else if err = insertPageDB(ctx, s.PageCollection, page); err != nil {
		return nil, err
	} else {
		return page, err
	}
}

func (s Server) GetPage(ctx context.Context, req *pb.GetPageReq) (*pb.Page, error) {
	return getPageDB(ctx, s.PageCollection, req.PageId)
}

func (s Server) GetAllPages(_ *empty.Empty, stream pb.PageService_GetAllPagesServer) (err error) {
	if cur, err := getAllDB(stream.Context(), s.PageCollection); err != nil {
		return err
	} else {

		var page *pb.Page
		for cur.Next(stream.Context()) {
			if err = cur.Decode(&page); err != nil {
				return err
			} else if err = stream.Send(page); err != nil {
				return err
			}
		}
		if err = cur.Close(stream.Context()); err != nil {
			return err
		}
		stream.Context().Done()
		return nil

	}
}

func (s Server) UpdatePage(ctx context.Context, req *pb.UpdatePageReq) (*pb.Page, error) {
	var err error
	if err = pageValidator(ctx, s.TilesCollection, req.Page); err != nil {
		return nil, err
	} else if err = updatePageDB(ctx, s.PageCollection, req.PageId, req.Page); err != nil {
		return nil, err
	} else {
		return req.Page, nil
	}
}

func (s Server) DeletePage(ctx context.Context, req *pb.DeletePageReq) (*empty.Empty, error) {
	return &empty.Empty{}, deletePageDB(ctx, s.PageCollection, req.PageId)
}

func (s Server) ResolvePage(ctx context.Context, req *pb.GetPageReq) (*pb.ResultPage, error) {
	// get Page from id
	if page, err := s.GetPage(ctx, req); err != nil {
		return nil, err
	} else {
		return pageResolver(ctx, s.TilesCollection, page)
	}
}

func (s Server) Catalog(_ *empty.Empty, stream pb.PageService_CatalogServer) error {
	if cur, err := getAllDB(stream.Context(), s.PageCollection); err != nil {
		return err
	} else {
		var page *pb.Page
		for cur.Next(stream.Context()) {
			if err = cur.Decode(&page); err != nil {
				return err
			} else if resultPage, err := pageResolver(stream.Context(), s.TilesCollection, page); err != nil {
				return err
			} else if err = stream.Send(resultPage); err != nil {
				return err
			}
		}
		if err = cur.Close(stream.Context()); err != nil {
			return err
		}
		stream.Context().Done()
		return nil
	}
}

func (s Server) GetDropDown(ctx context.Context, req *pb.DropDownReq) (*pb.DropDownResp, error) {
	return distinctValueDB(ctx, s.TilesCollection, req.Field)
}
