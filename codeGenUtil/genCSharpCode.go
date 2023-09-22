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

	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta("                            CSharp Code Gen Start                                  "))
	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta(""))

	var (
		ctx        = context.Background()
		in         gendao.CGenDaoInput
		db         gdb.DB
		tbNames    []string
		classNames []string
		start      = time.Now()
	)

	// 运行前检查
	checkBefore("cSharp")

	// 加载配置
	loadConfig(ctx, &in)

	// 获取数据库对象
	db = getDB(in)

	// 设置tables
	if len(strings.TrimSpace(in.Tables)) == 0 {
		tablesRes, err := db.Query(ctx, consts.SQL_SERVER_TABLES)
		if err != nil {
			log.Fatalf(color.Cyan("%s 获取业务表信息失败"), err)
		}
		in.Tables = tablesRes[0]["TableNames"].String()
		log.Printf(color.Magenta("配置信息 tables为空,默认获取当前数据库所有表: ")+color.Green("%s"), in.Tables)
		log.Printf("")
	}

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

	genCommon(in.ProjectName)
	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta("                        CSharp Code Gen End Cost :%s                                "), time.Since(start))
	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
}

func genCSharpBizCode(tab Table) {
	basePath := "./template/cSharp/"
	t2, _ := template.ParseFiles(basePath + "model.cs.template")
	t3, _ := template.ParseFiles(basePath + "controller.cs.template")
	t4, _ := template.ParseFiles(basePath + "interface.cs.template")
	t5, _ := template.ParseFiles(basePath + "processor.cs.template")

	genByTemplate(t2, "./internal/cSharp/model/"+tab.OriginalTableName+".cs", tab)
	genByTemplate(t3, "./internal/cSharp/controller/"+tab.TableName+"Controller.cs", tab)
	genByTemplate(t4, "./internal/cSharp/interface/I"+tab.TableName+"Processor.cs", tab)
	genByTemplate(t5, "./internal/cSharp/processor/"+tab.TableName+"Processor.cs", tab)
}

func genByTemplate(t *template.Template, filename string, tab Table) {
	if t == nil {
		log.Printf(color.Red("模板文件缺失,导致 %s 生成失败"), filename)
		return
	}

	var b bytes.Buffer
	t.Execute(&b, tab)
	fileCreate(b, filename)
}

func genAddScoped(tbNames []string) {
	var imports []string
	imports = append(imports, "// Repo 注入")
	for i := range tbNames {
		im := "builder.Services.AddScoped<" + tbNames[i] + "Repo>();"
		imports = append(imports, im)
	}

	imports = append(imports, "")
	imports = append(imports, "// Processor 注入")
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
	log.Printf(color.Green("program.cs 生成完毕"))
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
	tab := Table{Imports: imports, ProjectName: tabs[0].ProjectName}
	t1, _ := template.ParseFiles(basePath + "repo.cs.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, "./internal/cSharp/repo/repo.cs")
	log.Printf(color.Green("repo.cs 生成完毕 "))
}

func genCommon(projectName string) {
	basePath := "./template/cSharp/"
	t2, _ := template.ParseFiles(basePath + "contextUtils.cs.template")
	tab := Table{ProjectName: projectName}
	genByTemplate(t2, "./internal/cSharp/common/"+"ContextUtils.cs", tab)

	t1, _ := template.ParseFiles(basePath + "r_result.cs.template")
	genByTemplate(t1, "./internal/cSharp/common/"+"R_Result.cs", tab)
}
