/**
  @author: 志强
  @since: 2023/7/24
  @desc: 代码生成工具类-Golang
**/

package codeGenUtil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/internal/consts"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
)

func GenGolangCode() {
	var (
		err        error
		ctx        = context.Background()
		in         gendao.CGenDaoInput
		db         gdb.DB
		tbNames    []string
		classNames []string
	)

	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("./hack/config.yaml")
	in = gendao.CGenDaoInput{
		Path:       "internal",
		EntityPath: "model/entity",
		DoPath:     "model/do",
		DaoPath:    "dao",
	}

	err = g.Cfg().MustGet(
		ctx,
		fmt.Sprintf(`%s.%d`, gendao.CGenDaoConfig, 0),
	).Scan(&in)
	if err != nil {
		log.Fatalf(`invalid configuration of "%s": %+v`, gendao.CGenDaoConfig, err)
	}

	db = getDB(in)

	// 生成dao层代码
	err = gendao.Dao(context.Background(), in, db)
	if err != nil {
		return
	}

	fieldMap, err := db.Query(ctx, consts.SQL_SERVER_ShowTableStatus)
	tabs := make([]Table, 0)
	for i := range fieldMap {
		oriTbName := fieldMap[i]["Name"].String()
		tbName := strings.Replace(oriTbName, "T_", "", 1)
		if !strings.Contains(in.Tables, tbName) {
			continue
		}

		newTbName := gstr.CaseCamel(tbName)
		tab := Table{
			ClassName:         newTbName,
			TableName:         tbName,
			TableComment:      fieldMap[i]["Comment"].String(),
			OriginalTableName: oriTbName,
			Path:              toLow(newTbName),
		}
		tabs = append(tabs, tab)
		tbNames = append(tbNames, tab.TableName)
		classNames = append(classNames, tab.ClassName)
	}

	// 生成业务代码
	for i := range tabs {
		tb := tabs[i]
		definition := genBaseDefinition(db, tb.OriginalTableName, "Golang")
		tb.BaseDefinition = definition
		genBizCode(tb)
	}

	log.Printf(color.Cyan("%s 业务代码生成完毕"), tbNames)

	// 生成logic/logic.go
	genLogicImport(tbNames)

	// 生成业务路由注册 router/bizRouter.go
	genBizRouter(classNames)
}

func genBizCode(tab Table) {
	tab.ProjectName = getProjectName()
	basePath := "./template/golang/"
	t1, _ := template.ParseFiles(basePath + "api.go.template")
	t2, _ := template.ParseFiles(basePath + "model.go.template")
	t3, _ := template.ParseFiles(basePath + "controller.go.template")
	t4, _ := template.ParseFiles(basePath + "logic.go.template")
	t5, _ := template.ParseFiles(basePath + "service.go.template")

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

func genLogicImport(tbNames []string) {
	var pName = getProjectName()

	var imports []string

	for i := range tbNames {
		im := "\"" + pName + "/internal/logic/" + tbNames[i] + "\""
		imports = append(imports, im)
	}

	tab := Table{Imports: imports, ProjectName: getProjectName()}

	basePath := "./template/golang/"
	t1, _ := template.ParseFiles(basePath + "logic_all.go.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/logic/logic.go")
	log.Printf(color.Cyan("logic.go 生成完毕 引入包汇总: %s"), tbNames)
}

func genBizRouter(classNames []string) {
	tab := BizRouter{ClassNames: classNames, ProjectName: getProjectName()}

	basePath := "./template/golang/"
	t1, _ := template.ParseFiles(basePath + "bizRouter.go.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/router/bizRouter.go")
	log.Printf(color.Cyan("bizRouter.go 生成完毕 引入包汇总: %s"), classNames)
}
