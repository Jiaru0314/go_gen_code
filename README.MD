## 基于GoFrame v2 版本的代码生成工具

### 使用方法

#### Step1 下载代码生成工具库
`go get -u github.com/Jiaru0314/go_gen_code`

###  Step2
    修改  `/hack/config.yaml` 文件中数据库连接信息 以及配置需要生成的表 可参考GoFrame官方文档

```yaml
gfcli:
  gen:
    dao:
      - link: "mysql:root:123456@tcp(127.0.0.1:3306)/your_database"
        tables: "document"
        descriptionTag: true
        noModelComment: true
```

### Step3
    执行/main.go 中的 gencode.GenALl()方法生成dao层以及业务代码

```go
import (
	_ "context"
	"github.com/Jiaru0314/go_gen_code/codeGenUtil"
	_ "github.com/Jiaru0314/go_gen_code/gendao"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	codeGenUtil.GenALl()
}
```

    生成结果示例:
```
2023/07/26 10:30:01 generated: internal\dao\document.go
2023/07/26 10:30:01 generated: internal\dao\internal\document.go
2023/07/26 10:30:01 generated: internal\model/do\document.go
2023/07/26 10:30:01 generated: internal\model/entity\document.go
[document] Dao 层代码生成完毕

2023/07/26 10:30:01 generated: ./api/document.go
2023/07/26 10:30:01 generated: ./internal/model/document.go
2023/07/26 10:30:01 generated: ./internal/controller/document.go
2023/07/26 10:30:01 generated: ./internal/logic/document/document.go
2023/07/26 10:30:01 generated: ./internal/service/document.go

2023/07/26 10:30:01 [document] 业务代码生成完毕

```


