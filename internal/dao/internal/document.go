// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocumentDao is the data access object for table document.
type DocumentDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns DocumentColumns // columns contains all the column names of Table for convenient usage.
}

// DocumentColumns defines and stores column names for table document.
type DocumentColumns struct {
	Id        string // 编码
	Title     string // 标题
	Author    string // 作者
	Content   string // 内容
	Path      string // 文档路径
	PublishAt string // 发布时间
	CreatedAt string // 创建时间
	UpdatedAt string // 修改时间
	DeletedAt string // 删除时间
	CreateBy  string // 创建人
	UpdateBy  string // 修改人
	Cost      string // 处理耗时
}

// documentColumns holds the columns for table document.
var documentColumns = DocumentColumns{
	Id:        "id",
	Title:     "title",
	Author:    "author",
	Content:   "content",
	Path:      "path",
	PublishAt: "publish_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
	CreateBy:  "create_by",
	UpdateBy:  "update_by",
	Cost:      "cost",
}

// NewDocumentDao creates and returns a new DAO object for table data access.
func NewDocumentDao() *DocumentDao {
	return &DocumentDao{
		group:   "",
		table:   "document",
		columns: documentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DocumentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DocumentDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DocumentDao) Columns() DocumentColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DocumentDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DocumentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DocumentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
