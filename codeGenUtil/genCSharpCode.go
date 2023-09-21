/**
  @author: 志强
  @since: 2023/7/24
  @desc: 代码生成工具类-Golang
**/

package codeGenUtil

import (
	"bytes"
	"context"
	"log"
	_ "path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/internal/consts"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
)

func GenCSharpCode() {
	var (
		ctx        = context.Background()
		in         gendao.CGenDaoInput
		db         gdb.DB
		tbNames    []string
		classNames []string
		start      = time.Now()
	)

	// 加载配置
	loadConfig(ctx, &in)

	// 获取数据库对象
	db = getDB(in)

	fieldMap, err := db.Query(ctx, consts.SQL_SERVER_ShowTableStatus)
	if err != nil {
		log.Fatalf(color.Cyan("%s 获取业务表信息失败"), err)
	}
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
			ProjectName:       in.ProjectName,
		}
		tabs = append(tabs, tab)
		tbNames = append(tbNames, tab.TableName)
		classNames = append(classNames, tab.ClassName)
	}

	// 生成业务代码
	for i := range tabs {
		tb := tabs[i]
		tb.BaseDefinition = genBaseDefinition(db, tb.OriginalTableName, "CSharp")
		genCSharpBizCode(tb)
	}

	genAddScoped(tbNames)

	genRepo(tabs)

	log.Printf(color.Cyan("%s 业务代码生成完毕, 耗时：%s (ms)"), tbNames, time.Since(start))
}

func genCSharpBizCode(tab Table) {
	basePath := "./template/cSharp/"
	t2, _ := template.ParseFiles(basePath + "model.cs.template")
	t3, _ := template.ParseFiles(basePath + "controller.cs.template")
	t4, _ := template.ParseFiles(basePath + "interface.cs.template")
	t5, _ := template.ParseFiles(basePath + "processor.cs.template")

	var b2, b3, b4, b5 bytes.Buffer
	t2.Execute(&b2, tab)
	t3.Execute(&b3, tab)
	t4.Execute(&b4, tab)
	t5.Execute(&b5, tab)

	fileCreate(b2, "./internal/cSharp/model/"+tab.OriginalTableName+".cs")
	fileCreate(b3, "./internal/cSharp/controller/"+tab.TableName+"Controller.cs")
	fileCreate(b4, "./internal/cSharp/interface/"+"I"+tab.TableName+"Processor.cs")
	fileCreate(b5, "./internal/cSharp/processor/"+tab.TableName+"Processor.cs")
}

func genAddScoped(tbNames []string) {
	var imports []string

	for i := range tbNames {
		im := "builder.Services.AddScoped<I" + tbNames[i] + "Processor, " + tbNames[i] + "Processor>();"
		imports = append(imports, im)
	}

	basePath := "./template/cSharp/"
	tab := Table{Imports: imports}
	t1, _ := template.ParseFiles(basePath + "program.cs.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/cSharp/program.cs")
	log.Printf(color.Cyan("program.cs 生成完毕 引入interface汇总: %s"), tbNames)
}

func genRepo(tabs []Table) {
	var imports []string

	for i := range tabs {
		tab := tabs[i]
		imports = append(imports, "    /// <summary>")
		imports = append(imports, "    /// "+tab.TableComment)
		imports = append(imports, "    /// </summary>")

		im := "    public class " + tab.TableName + "Repo : BaseRepository<" + tab.OriginalTableName + ">"
		imports = append(imports, im)
		imports = append(imports, "    {")
		imports = append(imports, "    }")
		imports = append(imports, "")

	}

	basePath := "./template/cSharp/"
	tab := Table{Imports: imports}
	t1, _ := template.ParseFiles(basePath + "repo.cs.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/cSharp/repo/repo.cs")
	log.Printf(color.Cyan("repo.cs 生成完毕 引入interface汇总"))
}
