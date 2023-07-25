package model

// DocumentBase Document基类
type DocumentBase struct {
}

// DocumentInfo Document全部内容
type DocumentInfo struct {
    DocumentBase
    Id int
}

// AddDocumentInput 创建文档请求内容
type AddDocumentInput struct {
	DocumentBase
}

// AddDocumentOutput 创建文档返回响应
type AddDocumentOutput struct {
	Id int `json:"id"`
}

// DeleteDocumentInput 删除文档请求内容
type DeleteDocumentInput struct {
	Id int
}

// DeleteDocumentOutPut 删除文档响应
type DeleteDocumentOutPut struct {
	Success bool
}

// UpdateDocumentInput 修改文档请求内容
type UpdateDocumentInput struct {
	DocumentBase
	Id int
}

// UpdateDocumentOutPut 修改文档响应
type UpdateDocumentOutPut struct {
	Success bool
}

// PageDocumentInput 分页获取文档列表请求内容
type PageDocumentInput struct {
	Page int // 分页号码
	Size int // 分页数量，最大50
	Sort int // 排序类型
	DocumentBase
}

// PageDocumentOutput 分页获取文档列表请求响应
type PageDocumentOutput struct {
	List  []DocumentInfo    `json:"list"`  // 列表
	Page  int                     `json:"page"`  // 分页码
	Size  int                     `json:"size"`  // 分页数量
	Total int                     `json:"total"` // 数据总数
}

// DetailDocumentInput 根据ID查询文档 详情请求内容
type DetailDocumentInput struct {
	Id int
}

// DetailDocumentOutput 查询文档 详情响应
type DetailDocumentOutput struct {
	DocumentInfo
}


