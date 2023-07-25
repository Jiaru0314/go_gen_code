package Document

import (
	"context"
	"go_gen_code/internal/dao"
	"go_gen_code/internal/model"
	"go_gen_code/internal/service"
)

type sDocument struct{}

func init() {
	service.RegisterDocument(&sDocument{})
}

func (*sDocument) Add(ctx context.Context, in model.AddDocumentInput) (out *model.AddDocumentOutput, err error) {
	id, err := dao.Document.Ctx(ctx).Data(in.DocumentBase).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &model.AddDocumentOutput{Id: int(id)}, err
}

func (*sDocument) Update(ctx context.Context, in model.UpdateDocumentInput) (err error) {
	if _, err = dao.Document.Ctx(ctx).Data(in).FieldsEx(in.Id).Where(dao.Document.Columns().Id, in.Id).Update(); err != nil {
		return err
	}
	return nil
}

func (*sDocument) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.Document.Ctx(ctx).Where(dao.Document.Columns().Id, id).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (*sDocument) Page(ctx context.Context, in model.PageDocumentInput) (out *model.PageDocumentOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.Document.Ctx(ctx)
	//2. 实例化响应结构体
	out = &model.PageDocumentOutput{}
	out.Page, out.Size = in.Page, in.Size
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.DocumentInfo, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	var list []model.DocumentInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	out.List = list
	return
}
