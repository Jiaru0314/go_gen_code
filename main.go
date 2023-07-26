package main

import (
	_ "context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/Jiaru0314/go_gen_code/codeGenUtil"
	_ "github.com/Jiaru0314/go_gen_code/gendao"
)

func main() {
	codeGenUtil.GenALl()
}
