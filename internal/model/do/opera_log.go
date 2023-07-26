// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperaLog is the golang structure of table sys_opera_log for DAO operations like Where/Data.
type OperaLog struct {
	g.Meta        `orm:"table:sys_opera_log, do:true"`
	Id            interface{} // 主键编码
	Title         interface{} // 操作模块
	BusinessType  interface{} // 操作类型
	BusinessTypes interface{} // BusinessTypes
	Method        interface{} // 函数
	RequestMethod interface{} // 请求方式: GET POST PUT DELETE
	OperatorType  interface{} // 操作类型
	OperName      interface{} // 操作者
	DeptName      interface{} // 部门名称
	OperUrl       interface{} // 访问地址
	OperIp        interface{} // 客户端ip
	OperLocation  interface{} // 访问位置
	OperParam     interface{} // 请求参数
	Status        interface{} // 操作状态 1:正常 2:关闭
	OperTime      *gtime.Time // 操作时间
	JsonResult    interface{} // 返回数据
	Remark        interface{} // 备注
	LatencyTime   interface{} // 耗时
	UserAgent     interface{} // ua
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 最后更新时间
	CreateBy      interface{} // 创建者
	UpdateBy      interface{} // 更新者
}
