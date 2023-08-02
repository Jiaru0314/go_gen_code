package consts

var (
	MYSQL_ShowTableStatus      = "SHOW TABLE STATUS"
	SQL_SERVER_ShowTableStatus = "SELECT t.name AS Name,ep.name AS CommentName,ep.value AS Comment FROM sys.tables t LEFT JOIN sys.extended_properties ep ON ep.major_id = t.object_id AND ep.minor_id = 0;"
)
