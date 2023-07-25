package SysOperaLog

import (
	"context"
	"go_gen_code/internal/dao"
	"go_gen_code/internal/model"
	"go_gen_code/internal/service"
)

type sSysOperaLog struct{}

func init() {
	service.RegisterSysOperaLog(&sSysOperaLog{})
}

func (*sSysOperaLog) Add(ctx context.Context, in model.AddSysOperaLogInput) (out *model.AddSysOperaLogOutput, err error) {
	id, err := dao.SysOperaLog.Ctx(ctx).Data(in.SysOperaLogBase).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &model.AddSysOperaLogOutput{Id: int(id)}, err
}

func (*sSysOperaLog) Update(ctx context.Context, in model.UpdateSysOperaLogInput) (err error) {
	if _, err = dao.SysOperaLog.Ctx(ctx).Data(in).FieldsEx(in.Id).Where(dao.SysOperaLog.Columns().Id, in.Id).Update(); err != nil {
		return err
	}
	return nil
}

func (*sSysOperaLog) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.SysOperaLog.Ctx(ctx).Where(dao.SysOperaLog.Columns().Id, id).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (*sSysOperaLog) Page(ctx context.Context, in model.PageSysOperaLogInput) (out *model.PageSysOperaLogOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.SysOperaLog.Ctx(ctx)
	//2. 实例化响应结构体
	out = &model.PageSysOperaLogOutput{}
	out.Page, out.Size = in.Page, in.Size
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.SysOperaLogInfo, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	var list []model.SysOperaLogInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	out.List = list
	return
}
