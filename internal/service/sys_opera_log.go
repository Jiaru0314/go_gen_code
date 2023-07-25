package service

import (
	"context"
	"go_gen_code/internal/model"
)

type (
	ISysOperaLog interface {
		Add(ctx context.Context, in model.AddSysOperaLogInput) (out *model.AddSysOperaLogOutput, err error)
		Update(ctx context.Context, in model.UpdateSysOperaLogInput) (err error)
		Delete(ctx context.Context, id int) (err error)
		Page(ctx context.Context, in model.PageSysOperaLogInput) (out *model.PageSysOperaLogOutput, err error)
	}
)

var (
	localSysOperaLog ISysOperaLog
)

func SysOperaLog() ISysOperaLog {
	if localSysOperaLog == nil {
		panic("implement not found for interface ISysOperaLog, forgot register?")
	}
	return localSysOperaLog
}

func RegisterSysOperaLog(i ISysOperaLog) {
	localSysOperaLog = i
}
