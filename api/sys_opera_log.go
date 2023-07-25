package api

import "github.com/gogf/gf/v2/frame/g"

type SysOperaLogBase struct {
}

// AddSysOperaLogReq 新增
type AddSysOperaLogReq struct {
	g.Meta `path:"/sys_opera_log" method:"post" tags:"系统操作日志" summary:"新增"`
	SysOperaLogBase
}

type AddSysOperaLogRes struct {
	Id int `json:"id" dc:"id"`
}

// DeleteSysOperaLogReq 删除
type DeleteSysOperaLogReq struct {
	g.Meta `path:"/sys_opera_log" method:"delete" tags:"系统操作日志" summary:"删除"`
	Id     int `json:"id" v:"required#系统操作日志ID不能为空" dc:"id"`
}

type DeleteSysOperaLogRes struct{
	Success bool `json:"success"      description:"是否成功"`
}

// UpdateSysOperaLogReq 更新
type UpdateSysOperaLogReq struct {
	g.Meta `path:"/sys_opera_log" method:"put" tags:"系统操作日志" summary:"更新"`
	Id     int `json:"id" v:"required#系统操作日志ID不能为空" dc:"id"`
	SysOperaLogBase
}

type UpdateSysOperaLogRes struct{
	Success bool `json:"success"      description:"是否成功"`
}

// PageSysOperaLogReq 获取
type PageSysOperaLogReq struct {
	g.Meta `path:"/sys_opera_log/page" method:"get" tags:"系统操作日志" summary:"获取分页"`
	CommonPaginationReq
}

type PageSysOperaLogRes struct {
	CommonPaginationRes
}

