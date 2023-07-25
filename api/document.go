package api

import "github.com/gogf/gf/v2/frame/g"

type DocumentBase struct {
}

// AddDocumentReq 新增
type AddDocumentReq struct {
	g.Meta `path:"/document" method:"post" tags:"文档" summary:"新增"`
	DocumentBase
}

type AddDocumentRes struct {
	Id int `json:"id" dc:"id"`
}

// DeleteDocumentReq 删除
type DeleteDocumentReq struct {
	g.Meta `path:"/document" method:"delete" tags:"文档" summary:"删除"`
	Id     int `json:"id" v:"required#文档ID不能为空" dc:"id"`
}

type DeleteDocumentRes struct{
	Success bool `json:"success"      description:"是否成功"`
}

// UpdateDocumentReq 更新
type UpdateDocumentReq struct {
	g.Meta `path:"/document" method:"put" tags:"文档" summary:"更新"`
	Id     int `json:"id" v:"required#文档ID不能为空" dc:"id"`
	DocumentBase
}

type UpdateDocumentRes struct{
	Success bool `json:"success"      description:"是否成功"`
}

// PageDocumentReq 获取
type PageDocumentReq struct {
	g.Meta `path:"/document/page" method:"get" tags:"文档" summary:"获取分页"`
	CommonPaginationReq
}

type PageDocumentRes struct {
	CommonPaginationRes
}

