// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperaLog is the golang structure for table opera_log.
type OperaLog struct {
	Id            int64       `json:"id"             description:"主键编码"`
	Title         string      `json:"title"          description:"操作模块"`
	BusinessType  string      `json:"business_type"  description:"操作类型"`
	BusinessTypes string      `json:"business_types" description:"BusinessTypes"`
	Method        string      `json:"method"         description:"函数"`
	RequestMethod string      `json:"request_method" description:"请求方式: GET POST PUT DELETE"`
	OperatorType  string      `json:"operator_type"  description:"操作类型"`
	OperName      string      `json:"oper_name"      description:"操作者"`
	DeptName      string      `json:"dept_name"      description:"部门名称"`
	OperUrl       string      `json:"oper_url"       description:"访问地址"`
	OperIp        string      `json:"oper_ip"        description:"客户端ip"`
	OperLocation  string      `json:"oper_location"  description:"访问位置"`
	OperParam     string      `json:"oper_param"     description:"请求参数"`
	Status        string      `json:"status"         description:"操作状态 1:正常 2:关闭"`
	OperTime      *gtime.Time `json:"oper_time"      description:"操作时间"`
	JsonResult    string      `json:"json_result"    description:"返回数据"`
	Remark        string      `json:"remark"         description:"备注"`
	LatencyTime   string      `json:"latency_time"   description:"耗时"`
	UserAgent     string      `json:"user_agent"     description:"ua"`
	CreatedAt     *gtime.Time `json:"created_at"     description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updated_at"     description:"最后更新时间"`
	CreateBy      int64       `json:"create_by"      description:"创建者"`
	UpdateBy      int64       `json:"update_by"      description:"更新者"`
}
