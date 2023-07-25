// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Jiaru0314/go_gen_code/gendao/consts"
	"github.com/Jiaru0314/go_gen_code/gendao/utils"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func generateDao(ctx context.Context, in CGenDaoInternalInput) {
	var (
		dirPathDao         = gfile.Join(in.Path, in.DaoPath)
		dirPathDaoInternal = gfile.Join(dirPathDao, "internal")
	)
	if in.Clear {
		doClear(ctx, dirPathDao, true)
	}
	for i := 0; i < len(in.TableNames); i++ {
		generateDaoSingle(ctx, generateDaoSingleInput{
			CGenDaoInternalInput: in,
			TableName:            in.TableNames[i],
			NewTableName:         in.NewTableNames[i],
			DirPathDao:           dirPathDao,
			DirPathDaoInternal:   dirPathDaoInternal,
		})
	}
}

type generateDaoSingleInput struct {
	CGenDaoInternalInput
	TableName          string // TableName specifies the table name of the table.
	NewTableName       string // NewTableName specifies the prefix-stripped name of the table.
	DirPathDao         string
	DirPathDaoInternal string
}

// generateDaoSingle generates the dao and model content of given table.
func generateDaoSingle(ctx context.Context, in generateDaoSingleInput) {
	// Generating table data preparing.
	fieldMap, err := in.DB.TableFields(ctx, in.TableName)
	if err != nil {
		log.Fatalf(`fetching tables fields failed for table "%s": %+v`, in.TableName, err)
	}
	var (
		tableNameCamelCase      = gstr.CaseCamel(in.NewTableName)
		tableNameCamelLowerCase = gstr.CaseCamelLower(in.NewTableName)
		tableNameSnakeCase      = gstr.CaseSnake(in.NewTableName)
		importPrefix            = in.ImportPrefix
	)
	if importPrefix == "" {
		importPrefix = utils.GetImportPath(gfile.Join(in.Path, in.DaoPath))
	} else {
		importPrefix = gstr.Join(g.SliceStr{importPrefix, in.DaoPath}, "/")
	}

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}

	// dao - index
	generateDaoIndex(generateDaoIndexInput{
		generateDaoSingleInput:  in,
		TableNameCamelCase:      tableNameCamelCase,
		TableNameCamelLowerCase: tableNameCamelLowerCase,
		ImportPrefix:            importPrefix,
		FileName:                fileName,
	})

	// dao - internal
	generateDaoInternal(generateDaoInternalInput{
		generateDaoSingleInput:  in,
		TableNameCamelCase:      tableNameCamelCase,
		TableNameCamelLowerCase: tableNameCamelLowerCase,
		ImportPrefix:            importPrefix,
		FileName:                fileName,
		FieldMap:                fieldMap,
	})
}

type generateDaoIndexInput struct {
	generateDaoSingleInput
	TableNameCamelCase      string
	TableNameCamelLowerCase string
	ImportPrefix            string
	FileName                string
}

func generateDaoIndex(in generateDaoIndexInput) {
	path := gfile.Join(in.DirPathDao, in.FileName+".go")
	if in.OverwriteDao || !gfile.Exists(path) {
		indexContent := gstr.ReplaceByMap(
			getTemplateFromPathOrDefault(in.TplDaoIndexPath, consts.TemplateGenDaoIndexContent),
			g.MapStrStr{
				tplVarImportPrefix:            in.ImportPrefix,
				tplVarTableName:               in.TableName,
				tplVarTableNameCamelCase:      in.TableNameCamelCase,
				tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
			})
		indexContent = replaceDefaultVar(in.CGenDaoInternalInput, indexContent)
		if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
			log.Fatalf("writing content to '%s' failed: %v", path, err)
		} else {
			utils.GoFmt(path)
			log.Printf(color.Green("generated: %s"), path)
		}
	}
}

type generateDaoInternalInput struct {
	generateDaoSingleInput
	TableNameCamelCase      string
	TableNameCamelLowerCase string
	ImportPrefix            string
	FileName                string
	FieldMap                map[string]*gdb.TableField
}

func generateDaoInternal(in generateDaoInternalInput) {
	path := gfile.Join(in.DirPathDaoInternal, in.FileName+".go")
	modelContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoInternalPath, consts.TemplateGenDaoInternalContent),
		g.MapStrStr{
			tplVarImportPrefix:            in.ImportPrefix,
			tplVarTableName:               in.TableName,
			tplVarGroupName:               in.Group,
			tplVarTableNameCamelCase:      in.TableNameCamelCase,
			tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
			tplVarColumnDefine:            gstr.Trim(generateColumnDefinitionForDao(in.FieldMap)),
			tplVarColumnNames:             gstr.Trim(generateColumnNamesForDao(in.FieldMap)),
		})
	modelContent = replaceDefaultVar(in.CGenDaoInternalInput, modelContent)
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		log.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		log.Printf(color.Green("generated: %s"), path)
	}
}

// generateColumnNamesForDao generates and returns the column names assignment content of column struct
// for specified table.
func generateColumnNamesForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		field := fieldMap[name]
		array[index] = []string{
			"            #" + gstr.CaseCamel(field.Name) + ":",
			fmt.Sprintf(` #"%s",`, field.Name),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	namesContent := buffer.String()
	// Let's do this hack of table writer for indent!
	namesContent = gstr.Replace(namesContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(namesContent)
	return buffer.String()
}

// generateColumnDefinitionForDao generates and returns the column names definition for specified table.
func generateColumnDefinitionForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		var (
			field   = fieldMap[name]
			comment = gstr.Trim(gstr.ReplaceByArray(field.Comment, g.SliceStr{
				"\n", " ",
				"\r", " ",
			}))
		)
		array[index] = []string{
			"    #" + gstr.CaseCamel(field.Name),
			" # " + "string",
			" #" + fmt.Sprintf(`// %s`, comment),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	defineContent := buffer.String()
	// Let's do this hack of table writer for indent!
	defineContent = gstr.Replace(defineContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(defineContent)
	return buffer.String()
}
