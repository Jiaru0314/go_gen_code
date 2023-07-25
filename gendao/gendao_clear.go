// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"github.com/Jiaru0314/go_gen_code/gendao/utils"
	"log"

	"github.com/gogf/gf/v2/os/gfile"
)

func doClear(ctx context.Context, dirPath string, force bool) {
	files, err := gfile.ScanDirFile(dirPath, "*.go", true)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if force || utils.IsFileDoNotEdit(file) {
			if err = gfile.Remove(file); err != nil {
				log.Print(err)
			}
		}
	}
}
