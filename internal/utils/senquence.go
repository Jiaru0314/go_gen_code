package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func GetSequenceByTableName(tableName string) string {
	var tbName = strings.ReplaceAll(tableName, "t_", "SEQ_")
	var ctx = context.Background()
	var sql = fmt.Sprintf("SELECT SUBSTRING(CONVERT(VARCHAR,GETDATE(),112),3,6)+dbo.PadLeft(CONVERT(VARCHAR,NEXT VALUE FOR %s),'0',5)", tbName)
	query, _ := g.DB().Query(ctx, sql)
	m := query.List()[0][""]
	return gconv.String(m)
}
