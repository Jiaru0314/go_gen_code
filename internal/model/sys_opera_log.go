package model

// SysOperaLogBase SysOperaLog基类
type SysOperaLogBase struct {
}

// SysOperaLogInfo SysOperaLog全部内容
type SysOperaLogInfo struct {
    SysOperaLogBase
    Id int
}

// AddSysOperaLogInput 创建系统操作日志请求内容
type AddSysOperaLogInput struct {
	SysOperaLogBase
}

// AddSysOperaLogOutput 创建系统操作日志返回响应
type AddSysOperaLogOutput struct {
	Id int `json:"id"`
}

// DeleteSysOperaLogInput 删除系统操作日志请求内容
type DeleteSysOperaLogInput struct {
	Id int
}

// DeleteSysOperaLogOutPut 删除系统操作日志响应
type DeleteSysOperaLogOutPut struct {
	Success bool
}

// UpdateSysOperaLogInput 修改系统操作日志请求内容
type UpdateSysOperaLogInput struct {
	SysOperaLogBase
	Id int
}

// UpdateSysOperaLogOutPut 修改系统操作日志响应
type UpdateSysOperaLogOutPut struct {
	Success bool
}

// PageSysOperaLogInput 分页获取系统操作日志列表请求内容
type PageSysOperaLogInput struct {
	Page int // 分页号码
	Size int // 分页数量，最大50
	Sort int // 排序类型
	SysOperaLogBase
}

// PageSysOperaLogOutput 分页获取系统操作日志列表请求响应
type PageSysOperaLogOutput struct {
	List  []SysOperaLogInfo    `json:"list"`  // 列表
	Page  int                     `json:"page"`  // 分页码
	Size  int                     `json:"size"`  // 分页数量
	Total int                     `json:"total"` // 数据总数
}

// DetailSysOperaLogInput 根据ID查询系统操作日志 详情请求内容
type DetailSysOperaLogInput struct {
	Id int
}

// DetailSysOperaLogOutput 查询系统操作日志 详情响应
type DetailSysOperaLogOutput struct {
	SysOperaLogInfo
}


