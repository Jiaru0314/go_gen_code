// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OperaLogDao is the data access object for table sys_opera_log.
type OperaLogDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns OperaLogColumns // columns contains all the column names of Table for convenient usage.
}

// OperaLogColumns defines and stores column names for table sys_opera_log.
type OperaLogColumns struct {
	Id            string // 主键编码
	Title         string // 操作模块
	BusinessType  string // 操作类型
	BusinessTypes string // BusinessTypes
	Method        string // 函数
	RequestMethod string // 请求方式: GET POST PUT DELETE
	OperatorType  string // 操作类型
	OperName      string // 操作者
	DeptName      string // 部门名称
	OperUrl       string // 访问地址
	OperIp        string // 客户端ip
	OperLocation  string // 访问位置
	OperParam     string // 请求参数
	Status        string // 操作状态 1:正常 2:关闭
	OperTime      string // 操作时间
	JsonResult    string // 返回数据
	Remark        string // 备注
	LatencyTime   string // 耗时
	UserAgent     string // ua
	CreatedAt     string // 创建时间
	UpdatedAt     string // 最后更新时间
	CreateBy      string // 创建者
	UpdateBy      string // 更新者
}

// operaLogColumns holds the columns for table sys_opera_log.
var operaLogColumns = OperaLogColumns{
	Id:            "id",
	Title:         "title",
	BusinessType:  "business_type",
	BusinessTypes: "business_types",
	Method:        "method",
	RequestMethod: "request_method",
	OperatorType:  "operator_type",
	OperName:      "oper_name",
	DeptName:      "dept_name",
	OperUrl:       "oper_url",
	OperIp:        "oper_ip",
	OperLocation:  "oper_location",
	OperParam:     "oper_param",
	Status:        "status",
	OperTime:      "oper_time",
	JsonResult:    "json_result",
	Remark:        "remark",
	LatencyTime:   "latency_time",
	UserAgent:     "user_agent",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	CreateBy:      "create_by",
	UpdateBy:      "update_by",
}

// NewOperaLogDao creates and returns a new DAO object for table data access.
func NewOperaLogDao() *OperaLogDao {
	return &OperaLogDao{
		group:   "",
		table:   "sys_opera_log",
		columns: operaLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OperaLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OperaLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OperaLogDao) Columns() OperaLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OperaLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OperaLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OperaLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
