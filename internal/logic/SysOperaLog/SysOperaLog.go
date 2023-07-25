package sys_opera_log

import (
	"context"
	"go_gen_code/internal/dao"
	"go_gen_code/internal/model"
	"go_gen_code/internal/service"
)

type ssys_opera_log struct{}

func init() {
	service.Registersys_opera_log(&ssys_opera_log{})
}

func (*ssys_opera_log) Add(ctx context.Context, in model.Addsys_opera_logInput) (out *model.Addsys_opera_logOutput, err error) {
	id, err := dao.sys_opera_log.Ctx(ctx).Data(in.sys_opera_logBase).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &model.Addsys_opera_logOutput{Id: int(id)}, err
}

func (*ssys_opera_log) Update(ctx context.Context, in model.Updatesys_opera_logInput) (err error) {
	if _, err = dao.sys_opera_log.Ctx(ctx).Data(in).FieldsEx(in.Id).Where(dao.sys_opera_log.Columns().Id, in.Id).Update(); err != nil {
		return err
	}
	return nil
}

func (*ssys_opera_log) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.sys_opera_log.Ctx(ctx).Where(dao.sys_opera_log.Columns().Id, id).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (*ssys_opera_log) Page(ctx context.Context, in model.Pagesys_opera_logInput) (out *model.Pagesys_opera_logOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.sys_opera_log.Ctx(ctx)
	//2. 实例化响应结构体
	out = &model.Pagesys_opera_logOutput{}
	out.Page, out.Size = in.Page, in.Size
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.sys_opera_logInfo, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	var list []model.sys_opera_logInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	out.List = list
	return
}
