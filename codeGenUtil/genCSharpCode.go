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
	"github.com/olekukonko/tablewriter"

	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/internal/consts"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
)

// CSharp代码生成模板路径
const cSharpTmplPath = "./template/cSharp/"

func GenCSharpCode() {

	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta("                            CSharp Code Gen Start                                  "))
	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta(""))

	var (
		ctx     = context.Background()
		in      gendao.CGenDaoInput
		db      *gdb.DB
		tbNames []string
		start   = time.Now()
	)

	// 运行前检查
	checkBefore("cSharp")

	// 加载配置
	loadConfig(ctx, &in)

	// 获取数据库对象
	db = getDB(&in)

	// 设置tables
	if len(strings.TrimSpace(in.Tables)) == 0 {
		in.Tables = getTables(ctx, *db)
	}

	// 获取实体定义
	tabs := getTableDefinition(ctx, *db, &tbNames, in.Tables, in.ProjectName)

	// 生成业务代码
	for i := range tabs {
		tb := tabs[i]
		tb.BaseDefinition = genBaseDefinition(*db, tb.OriginalTableName, "CSharp")
		genCSharpBizCode(tb, in.OutPath)
	}

	genAddScoped(tbNames, in.OutPath, in.ProjectName)

	genRepo(tabs, in.OutPath)

	genCommon(in.ProjectName, in.OutPath)

	genMiddleware(in.ProjectName, in.OutPath)

	genMapper(tabs, in.OutPath)

	genAuditLogMap(tabs, in.OutPath)

	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	log.Printf(color.Magenta("                        CSharp Code Gen End Cost :%s                                "), time.Since(start))
	log.Printf(color.Magenta("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
}

func getTableDefinition(ctx context.Context, db gdb.DB, tbNames *[]string, tables, projectName string) []Table {
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
		if newTbName == "" {
			continue
		}
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
		*tbNames = append(*tbNames, tab.TableName)
	}
	return tabs
}

func getTables(ctx context.Context, db gdb.DB) string {
	tablesRes, err := db.Query(ctx, consts.SQL_SERVER_TABLES)
	if err != nil {
		log.Fatalf(color.Cyan("%s 获取业务表信息失败"), err)
	}
	res := tablesRes[0]["TableNames"].String()
	log.Printf(color.Magenta("配置信息 tables为空,默认获取当前数据库所有表: ")+color.Green("%s"), res)
	log.Printf("")
	return res
}

func genCSharpBizCode(tab Table, outPath string) {
	// 设置是否包含 IsEnable/ Status 这类字段 特殊处理
	if strings.Contains(tab.BaseDefinition, "IsEnable") {
		tab.HasEnable = true
	} else if strings.Contains(tab.BaseDefinition, "Status") {
		tab.HasStatus = true
	}

	t2, _ := template.ParseFiles(cSharpTmplPath + "model.cs.template")
	t3, _ := template.ParseFiles(cSharpTmplPath + "controller.cs.template")
	t4, _ := template.ParseFiles(cSharpTmplPath + "interface.cs.template")
	t5, _ := template.ParseFiles(cSharpTmplPath + "processor.cs.template")

	genByTemplate(t2, outPath+"Models/"+tab.OriginalTableName+".cs", tab)
	genByTemplate(t3, outPath+"API/Controllers/"+tab.TableName+"Controller.cs", tab)
	genByTemplate(t4, outPath+"Interface/I"+tab.TableName+"Processor.cs", tab)
	genByTemplate(t5, outPath+"Processor/"+tab.TableName+"Processor.cs", tab)
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

func genAddScoped(tbNames []string, outPath, projectName string) {
	var imports []string
	var middlewares []string

	imports = append(imports, "builder.Services.AddHttpContextAccessor();")
	imports = append(imports, "")

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

	middlewares = append(middlewares, "// 中间件注入")
	middlewares = append(middlewares, "app.UseMiddleware<AuditTrailMiddleware>();")
	middlewares = append(middlewares, "app.UseMiddleware<ErrorHandlerMiddleware>();")

	imports = append(imports, "")
	imports = append(imports, "//AutoMapper 注入")
	imports = append(imports, "AutoMapper.IConfigurationProvider config = new MapperConfiguration(")
	imports = append(imports, "    cfg => { cfg.AddProfile<MappingProfile>(); }")
	imports = append(imports, ");")
	imports = append(imports, "builder.Services.AddSingleton(config);")
	imports = append(imports, "builder.Services.AddScoped<IMapper, Mapper>();")

	tab := Table{Imports: imports, ProjectName: projectName, Middlewares: middlewares}
	t1, _ := template.ParseFiles(cSharpTmplPath + "program.cs.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, outPath+"API"+"/program.cs")
	log.Printf(color.Green("program.cs 生成完毕"))
}

func genRepo(tabs []Table, outPath string) {
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

	tab := Table{Imports: imports, ProjectName: tabs[0].ProjectName}
	t1, _ := template.ParseFiles(cSharpTmplPath + "repo.cs.template")
	var b1 bytes.Buffer
	t1.Execute(&b1, tab)
	fileCreate(b1, outPath+"Modules"+"/ObjectRepository.cs.cs")
	log.Printf(color.Green("repo.cs 生成完毕 "))
}

func genCommon(projectName string, outPath string) {
	tab := Table{ProjectName: projectName}

	t1, _ := template.ParseFiles(cSharpTmplPath + "r_result.cs.template")
	genByTemplate(t1, outPath+"Models/"+"R_Result.cs", tab)

	t2, _ := template.ParseFiles(cSharpTmplPath + "contextUtils.cs.template")
	genByTemplate(t2, outPath+"Processor/"+"ContextUtils.cs", tab)

	t3, _ := template.ParseFiles(cSharpTmplPath + "b_model.cs.template")
	genByTemplate(t3, outPath+"Models/"+"BaseModel.cs", tab)
}

func genMiddleware(projectName string, outPath string) {
	t2, _ := template.ParseFiles(cSharpTmplPath + "auditTrailMiddleware.cs.template")
	tab := Table{ProjectName: projectName}
	genByTemplate(t2, outPath+"Middlewares/"+"AuditTrailMiddleware.cs", tab)

	t1, _ := template.ParseFiles(cSharpTmplPath + "errorMiddleware.cs.template")
	genByTemplate(t1, outPath+"Middlewares/"+"ErrorMiddleware.cs", tab)
}

func genMapper(tabs []Table, outPath string) {
	buffer := bytes.NewBuffer(nil)
	array := make([]string, 0)
	var space = "      "
	for _, tab := range tabs {
		array = append(array, space)
		array = append(array, space+"CreateMap<Add"+tab.TableName+"Req,"+tab.OriginalTableName+">();")
		array = append(array, space+"CreateMap<Update"+tab.TableName+"Req,"+tab.OriginalTableName+">();")
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("\n")
	tw.Append(array)
	tw.Render()
	stContent := buffer.String()
	buffer.Reset()
	buffer.WriteString(stContent)

	t2, _ := template.ParseFiles(cSharpTmplPath + "mapper.cs.template")
	tab := Table{MapperConfiguration: stContent, ProjectName: tabs[0].ProjectName}
	genByTemplate(t2, outPath+"Models/Mapper/"+"MappingProfile.cs", tab)
}

func genAuditLogMap(tabs []Table, outPath string) {
	methods := make([]string, 0)
	methods = append(methods, "Add")
	methods = append(methods, "Update")
	methods = append(methods, "Delete")

	logMaps := make([]string, 0)
	logMaps = append(logMaps, "")

	// 生成业务代码
	for i := range tabs {
		tb := tabs[i]
		for m := range methods {
			key := "/" + tb.TableName + "/" + methods[m] + tb.TableName
			var val = ""
			if methods[m] == "Add" {
				val = "新增" + tb.TableComment
			}
			if methods[m] == "Update" {
				val = "修改" + tb.TableComment
			}
			if methods[m] == "Delete" {
				val = "删除" + tb.TableComment
			}
			// fmt.Println(" { \"" + key + "\", \"" + val + "\" },")
			logMaps = append(logMaps, "        { \""+key+"\", \""+val+"\" },")
		}
	}

	buffer := bytes.NewBuffer(nil)
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("\n")
	tw.Append(logMaps)
	tw.Render()
	stContent := buffer.String()
	buffer.Reset()
	buffer.WriteString(stContent)

	tab := Table{ProjectName: tabs[0].ProjectName, LogMaps: stContent}

	t1, _ := template.ParseFiles(cSharpTmplPath + "const.cs.template")
	genByTemplate(t1, outPath+"Modules/"+"CommonConst.cs", tab)
}
