// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Jiaru0314/go_gen_code/gendao/consts"
	"github.com/Jiaru0314/go_gen_code/gendao/utils"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func generateDo(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathDo = gfile.Join(in.Path, in.DoPath)
	if in.Clear {
		doClear(ctx, dirPathDo, false)
	}
	in.NoJsonTag = true
	in.DescriptionTag = false
	in.NoModelComment = false
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			log.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}
		var (
			newTableName        = in.NewTableNames[i]
			doFilePath          = gfile.Join(dirPathDo, gstr.CaseSnake(newTableName)+".go")
			structDefinition, _ = GenerateStructDefinition(ctx, GenerateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           gstr.CaseCamel(newTableName),
				FieldMap:             fieldMap,
				IsDo:                 true,
			})
		)
		// replace all types to interface{}.
		structDefinition, _ = gregex.ReplaceStringFuncMatch(
			"([A-Z]\\w*?)\\s+([\\w\\*\\.]+?)\\s+(//)",
			structDefinition,
			func(match []string) string {
				// If the type is already a pointer/slice/map, it does nothing.
				if !gstr.HasPrefix(match[2], "*") && !gstr.HasPrefix(match[2], "[]") && !gstr.HasPrefix(match[2], "map") {
					return fmt.Sprintf(`%s interface{} %s`, match[1], match[3])
				}
				return match[0]
			},
		)
		modelContent := generateDoContent(
			ctx,
			in,
			tableName,
			gstr.CaseCamel(newTableName),
			structDefinition,
		)
		err = gfile.PutContents(doFilePath, strings.TrimSpace(modelContent))
		if err != nil {
			log.Fatalf(`writing content to "%s" failed: %v`, doFilePath, err)
		} else {
			utils.GoFmt(doFilePath)
			log.Printf(color.Green("generated: %s"), doFilePath)
		}
	}
}

func generateDoContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string,
) string {
	doContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoDoPath, consts.TemplateGenDaoDoContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(ctx, structDefine, true, nil),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		},
	)
	doContent = replaceDefaultVar(in, doContent)
	return doContent
}
