/**
  @author: 志强
  @since: 2023/9/14
  @desc: 代码生成 通用定义/工具
**/

package codeGenUtil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/Jiaru0314/go_gen_code/gendao"
	"github.com/Jiaru0314/go_gen_code/gendao/utils"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
)

type Table struct {
	ProjectName       string
	ClassName         string
	TableName         string
	TableComment      string
	Imports           []string
	BaseDefinition    string
	OriginalTableName string
	Path              string
}

type BizRouter struct {
	ProjectName string
	ClassNames  []string
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

	utils.GoFmt(name)
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

func genBaseDefinition(db gdb.DB, tableName, language string) string {
	fields, _ := db.TableFields(context.Background(), tableName)

	in := gendao.GenerateStructDefinitionInput{
		TableName:  tableName,
		StructName: gstr.CaseCamel(tableName),
		FieldMap:   fields,
		IsDo:       false,
	}
	in.DB = db
	if language == "Golang" {
		return gendao.GenerateBaseDefinition(context.Background(), in)
	} else {
		return gendao.GenerateBaseDefinitionForCSharp(context.Background(), in)
	}
}

func toLow(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}

func loadConfig(ctx context.Context, in *gendao.CGenDaoInput) {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("./hack/config.yaml")

	err := g.Cfg().MustGet(
		ctx,
		fmt.Sprintf(`%s.%d`, gendao.CGenDaoConfig, 0),
	).Scan(&in)
	if err != nil {
		log.Fatalf(`数据库信息配置错误 "%s": %+v`, gendao.CGenDaoConfig, err)
	}

	if len(strings.TrimSpace(in.Link)) == 0 {
		log.Fatalf("数据库信息配置错误 link 为空")
	}

	log.Printf(color.Magenta("当前配置信息 link: ")+color.Green("%s"), in.Link)
	log.Printf(color.Magenta("当前配置信息 tables: ")+color.Green("%s"), in.Tables)
	log.Printf(color.Magenta(""))
}

// checkBefore 前置校验
func checkBefore(language string) {
	// 1、校验模板是否存在
	stat, _ := gfile.Stat("./template")
	if stat == nil {
		log.Fatal("模板文件不存在,请检查当前文件路径下是否存在 template文件夹")
	}

	stat, _ = gfile.Stat("./template/" + language)
	if stat == nil {
		log.Fatal("cSharp 模板文件不存在,请检查当前文件路径下是否存在 template/cSharp 文件夹")
	}

	stat, _ = gfile.Stat("./hack/config.yaml")
	if stat == nil {
		log.Fatal("运行配置文件不存在,请检查当前文件路径下是否存在 hack/config.yaml 文件")
	}
}
