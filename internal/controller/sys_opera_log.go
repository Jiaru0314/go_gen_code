package controller

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"go_gen_code/api"
	"go_gen_code/internal/model"
	"go_gen_code/internal/service"
)

var SysOperaLog = cSysOperaLog{}

type cSysOperaLog struct{}

func (*cSysOperaLog) Add(ctx context.Context, req *api.AddSysOperaLogReq) (res *api.AddSysOperaLogRes, err error) {
	input := model.AddSysOperaLogInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	out, err := service.SysOperaLog().Add(ctx, input)
	if err != nil {
		return nil, err
	}

	return &api.AddSysOperaLogRes{Id: out.Id}, nil
}

func (*cSysOperaLog) Delete(ctx context.Context, req *api.DeleteSysOperaLogReq) (res *api.DeleteSysOperaLogRes, err error) {
	if err = service.SysOperaLog().Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	return &api.DeleteSysOperaLogRes{
    		Success: true,
    	}, nil
}

func (*cSysOperaLog) Update(ctx context.Context, req *api.UpdateSysOperaLogReq) (res *api.UpdateSysOperaLogRes, err error) {
	input := model.UpdateSysOperaLogInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	if err = service.SysOperaLog().Update(ctx, input); err != nil {

		return nil, err
	}

	return &api.UpdateSysOperaLogRes{
    		Success: true,
    	}, nil
}

func (*cSysOperaLog) Page(ctx context.Context, req *api.PageSysOperaLogReq) (res *api.PageSysOperaLogRes, err error) {
	input := model.PageSysOperaLogInput{}
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}

	out, err := service.SysOperaLog().Page(ctx, input)
	if err != nil {
		return nil, err
	}

	err = gconv.Struct(out, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
