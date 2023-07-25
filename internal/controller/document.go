package controller

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"go_gen_code/api"
	"go_gen_code/internal/model"
	"go_gen_code/internal/service"
)

var Document = cDocument{}

type cDocument struct{}

func (*cDocument) Add(ctx context.Context, req *api.AddDocumentReq) (res *api.AddDocumentRes, err error) {
	input := model.AddDocumentInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	out, err := service.Document().Add(ctx, input)
	if err != nil {
		return nil, err
	}

	return &api.AddDocumentRes{Id: out.Id}, nil
}

func (*cDocument) Delete(ctx context.Context, req *api.DeleteDocumentReq) (res *api.DeleteDocumentRes, err error) {
	if err = service.Document().Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	return &api.DeleteDocumentRes{
    		Success: true,
    	}, nil
}

func (*cDocument) Update(ctx context.Context, req *api.UpdateDocumentReq) (res *api.UpdateDocumentRes, err error) {
	input := model.UpdateDocumentInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	if err = service.Document().Update(ctx, input); err != nil {

		return nil, err
	}

	return &api.UpdateDocumentRes{
    		Success: true,
    	}, nil
}

func (*cDocument) Page(ctx context.Context, req *api.PageDocumentReq) (res *api.PageDocumentRes, err error) {
	input := model.PageDocumentInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	out, err := service.Document().Page(ctx, input)
	if err != nil {
		return nil, err
	}

	err = gconv.Struct(out, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
