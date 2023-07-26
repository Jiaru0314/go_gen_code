// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Document is the golang structure of table document for DAO operations like Where/Data.
type Document struct {
	g.Meta    `orm:"table:document, do:true"`
	Id        interface{} // 编码
	Title     interface{} // 标题
	Author    interface{} // 作者
	Content   interface{} // 内容
	Path      interface{} // 文档路径
	PublishAt *gtime.Time // 发布时间
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 修改时间
	DeletedAt *gtime.Time // 删除时间
	CreateBy  interface{} // 创建人
	UpdateBy  interface{} // 修改人
	Cost      interface{} // 处理耗时
}
