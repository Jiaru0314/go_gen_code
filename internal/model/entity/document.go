// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Document is the golang structure for table document.
type Document struct {
	Id        uint        `json:"id"         description:"编码"`
	Title     string      `json:"title"      description:"标题"`
	Author    string      `json:"author"     description:"作者"`
	Content   string      `json:"content"    description:"内容"`
	Path      string      `json:"path"       description:"文档路径"`
	PublishAt *gtime.Time `json:"publish_at" description:"发布时间"`
	CreatedAt *gtime.Time `json:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at" description:"修改时间"`
	DeletedAt *gtime.Time `json:"deleted_at" description:"删除时间"`
	CreateBy  uint        `json:"create_by"  description:"创建人"`
	UpdateBy  uint        `json:"update_by"  description:"修改人"`
	Cost      float64     `json:"cost"       description:"处理耗时"`
}
