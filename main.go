package main

import (
	_ "context"
	"github.com/Jiaru0314/go_gen_code/codeGenUtil"
	_ "github.com/Jiaru0314/go_gen_code/gendao"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	codeGenUtil.GenALl()
}
