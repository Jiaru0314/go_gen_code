// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"github.com/Jiaru0314/go_gen_code/gendao/consts"
	"github.com/Jiaru0314/go_gen_code/gendao/utils"
	"github.com/Jiaru0314/go_gen_code/internal/utils/color"
	"log"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func generateEntity(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathEntity = gfile.Join(in.Path, in.EntityPath)
	if in.Clear {
		doClear(ctx, dirPathEntity, false)
	}
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			log.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}

		var (
			newTableName                    = in.NewTableNames[i]
			entityFilePath                  = gfile.Join(dirPathEntity, gstr.CaseSnake(newTableName)+".go")
			structDefinition, appendImports = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           gstr.CaseCamel(newTableName),
				FieldMap:             fieldMap,
				IsDo:                 false,
			})
			entityContent = generateEntityContent(
				ctx,
				in,
				newTableName,
				gstr.CaseCamel(newTableName),
				structDefinition,
				appendImports,
			)
		)

		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			log.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			log.Printf(color.Green("generated: %s"), entityFilePath)
		}
	}
}

func generateEntityContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string, appendImports []string,
) string {
	entityContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoEntityPath, consts.TemplateGenDaoEntityContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(ctx, structDefine, false, appendImports),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		},
	)
	entityContent = replaceDefaultVar(in, entityContent)
	return entityContent
}
