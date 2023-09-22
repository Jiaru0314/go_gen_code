package consts

var (
	MYSQL_ShowTableStatus      = "SHOW TABLE STATUS"
	SQL_SERVER_ShowTableStatus = "SELECT t.name AS Name,ep.name AS CommentName,ep.value AS Comment FROM sys.tables t LEFT JOIN sys.extended_properties ep ON ep.major_id = t.object_id AND ep.minor_id = 0;"
	SQL_SERVER_TABLES          = "SELECT STUFF((SELECT ', ' + TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' FOR XML PATH('')),1,2,'') AS TableNames;"
)
