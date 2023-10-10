/**
  @author: 志强
  @since: 2023/7/24
  @desc: 代码生成工具类-Csharp
**/

package codeGenUtil

import (
	"bytes"
	"context"
	"log"
	"strings"
	"text/template"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/internal/consts"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
)

const golangBasePath = "./template/golang/"

func GenGolangCode() {
	var (
		err        error
		ctx        = context.Background()
		in         gendao.CGenDaoInput
		db         *gdb.DB
		tbNames    []*string
		classNames []*string
	)

	in = gendao.CGenDaoInput{
		Path:       "internal",
		EntityPath: "model/entity",
		DoPath:     "model/do",
		DaoPath:    "dao",
	}

	// 加载配置
	loadConfig(ctx, &in)

	db = getDB(&in)

	// 生成dao层代码
	err = gendao.Dao(context.Background(), in, *db)
	if err != nil {
		return
	}

	tabs := getGolangTableDefinition(ctx, *db, tbNames, classNames, in.Tables, in.ProjectName)

	// 生成业务代码
	for i := range tabs {
		tb := tabs[i]
		definition := genBaseDefinition(*db, tb.OriginalTableName, "Golang")
		tb.BaseDefinition = definition
		genBizCode(tb)
	}

	log.Printf(color.Cyan("%s 业务代码生成完毕"), tbNames)

	// 生成logic/logic.go
	genLogicImport(tbNames)

	// 生成业务路由注册 router/bizRouter.go
	genBizRouter(classNames)
}

func getGolangTableDefinition(ctx context.Context, db gdb.DB, tbNames []*string, classNames []*string, tables, projectName string) []Table {
	fieldMap, err := db.Query(ctx, consts.SQL_SERVER_ShowTableStatus)
	if err != nil {
		log.Fatalf(color.Cyan("%s 获取业务表信息失败"), err)
	}
	tabs := make([]Table, 0)
	for i := range fieldMap {
		oriTbName := fieldMap[i]["Name"].String()
		tbName := strings.Replace(oriTbName, "T_", "", 1)
		if !strings.Contains(tables, tbName) {
			continue
		}

		newTbName := gstr.CaseCamel(tbName)
		tableComment := fieldMap[i]["Comment"].String()
		if tableComment != "" && strings.Contains(tableComment, "表") {
			tableComment = strings.ReplaceAll(tableComment, "表", "")
		}

		tab := Table{
			ClassName:         newTbName,
			TableName:         tbName,
			TableComment:      tableComment,
			OriginalTableName: oriTbName,
			Path:              toLow(newTbName),
			ProjectName:       projectName,
		}
		tabs = append(tabs, tab)
		tbNames = append(tbNames, &tab.TableName)
		classNames = append(classNames, &tab.ClassName)
	}
	return tabs
}

func genBizCode(tab Table) {
	tab.ProjectName = getProjectName()
	t1, _ := template.ParseFiles(golangBasePath + "api.go.template")
	t2, _ := template.ParseFiles(golangBasePath + "model.go.template")
	t3, _ := template.ParseFiles(golangBasePath + "controller.go.template")
	t4, _ := template.ParseFiles(golangBasePath + "logic.go.template")
	t5, _ := template.ParseFiles(golangBasePath + "service.go.template")

	var b1, b2, b3, b4, b5 bytes.Buffer
	var err error
	err = t1.Execute(&b1, tab)
	err = t2.Execute(&b2, tab)
	err = t3.Execute(&b3, tab)
	err = t4.Execute(&b4, tab)
	err = t5.Execute(&b5, tab)

	if err != nil {
		log.Printf(color.Red("%s 业务代码生成失败"), err)
		return
	}

	fileCreate(b1, "./api/"+tab.TableName+".go")
	fileCreate(b2, "./internal/model/"+tab.TableName+".go")
	fileCreate(b3, "./internal/controller/"+tab.TableName+".go")
	fileCreate(b4, "./internal/logic/"+tab.TableName+"/"+tab.TableName+".go")
	fileCreate(b5, "./internal/service/"+tab.TableName+".go")
}

func genLogicImport(tbNames []*string) {
	var pName = getProjectName()

	var imports []string

	for i := range tbNames {
		im := "\"" + pName + "/internal/logic/" + *tbNames[i] + "\""
		imports = append(imports, im)
	}

	tab := Table{Imports: imports, ProjectName: getProjectName()}

	t1, _ := template.ParseFiles(golangBasePath + "logic_all.go.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/logic/logic.go")
	log.Printf(color.Cyan("logic.go 生成完毕 引入包汇总: %s"), tbNames)
}

func genBizRouter(classNames []*string) {
	tab := BizRouter{ClassNames: classNames, ProjectName: getProjectName()}

	t1, _ := template.ParseFiles(golangBasePath + "bizRouter.go.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/router/bizRouter.go")
	log.Printf(color.Cyan("bizRouter.go 生成完毕 引入包汇总: %s"), classNames)
}
