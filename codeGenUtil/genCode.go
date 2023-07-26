/**
  @author: 志强
  @since: 2023/7/24
  @desc: 代码生成工具类
**/

package codeGenUtil

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/internal/consts"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Table struct {
	ProjectName  string
	ClassName    string
	TableName    string
	TableComment string
}

func getProjectName() string {
	currentDir, _ := os.Getwd()
	split := strings.Split(currentDir, "\\")
	return split[len(split)-1]
}

func getDir(filePath string) string {
	dirPath := filepath.Dir(filePath)
	return dirPath
}

func GenALl() {
	var (
		err     error
		ctx     context.Context
		in      gendao.CGenDaoInput
		db      gdb.DB
		tbNames []string
	)

	ctx = context.Background()
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

	fieldMap, err := db.Query(ctx, consts.ShowTableStatus)
	tabs := make([]Table, 0)
	for i := range fieldMap {
		tbName := fieldMap[i]["Name"].String()
		if !strings.Contains(in.Tables, tbName) {
			continue
		}

		newTbName := gstr.CaseCamel(tbName)
		tab := Table{
			ClassName:    newTbName,
			TableName:    tbName,
			TableComment: fieldMap[i]["Comment"].String(),
		}
		tabs = append(tabs, tab)
	}

	//生成业务代码
	for i := range tabs {
		genBizCode(tabs[i])
	}

	for i := range tabs {
		tbNames = append(tbNames, tabs[i].TableName)
	}

	log.Printf(color.Cyan("%s 业务代码生成完毕"), tbNames)

}

func genBizCode(tab Table) {
	tab.ProjectName = getProjectName()
	basePath := "./template/"
	t1, _ := template.ParseFiles(basePath + "api.go.template")
	t2, _ := template.ParseFiles(basePath + "model.go.template")
	t3, _ := template.ParseFiles(basePath + "controller.go.template")
	t4, _ := template.ParseFiles(basePath + "logic.go.template")
	t5, _ := template.ParseFiles(basePath + "service.go.template")

	var b1, b2, b3, b4, b5 bytes.Buffer
	t1.Execute(&b1, tab)
	t2.Execute(&b2, tab)
	t3.Execute(&b3, tab)
	t4.Execute(&b4, tab)
	t5.Execute(&b5, tab)

	fileCreate(b1, "./api/"+tab.TableName+".go")
	fileCreate(b2, "./internal/model/"+tab.TableName+".go")
	fileCreate(b3, "./internal/controller/"+tab.TableName+".go")
	fileCreate(b4, "./internal/logic/"+tab.TableName+"/"+tab.TableName+".go")
	fileCreate(b5, "./internal/service/"+tab.TableName+".go")
}

func fileCreate(content bytes.Buffer, name string) {

	dir := getDir(name)
	_, err := os.Stat(dir)
	if err != nil {
		// 文件夹不存在，创建
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
			return
		}
	}

	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Println(err)
	}
	file.Close()
	log.Printf(color.Green("generated: %s"), name)
}

func getDB(in gendao.CGenDaoInput) (res gdb.DB) {
	var (
		db  gdb.DB
		err error
	)
	// It uses user passed database configuration.
	var tempGroup = gtime.TimestampNanoStr()
	gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
		Link: in.Link,
	})
	if db, err = gdb.Instance(tempGroup); err != nil {
		log.Fatalf(`database initialization failed: %+v`, err)
	}

	if db == nil {
		log.Fatal(`database initialization failed, may be invalid database configuration`)
	}

	return db
}
