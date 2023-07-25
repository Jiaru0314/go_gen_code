package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"go_gen_code/internal/consts"
	"go_gen_code/internal/controller"
)

var (
	Main = gcmd.Command{
		Name:  consts.ProjectName,
		Usage: consts.ProjectUsage,
		Brief: consts.ProjectUsage,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			registerBizRouter(s)
			s.Run()
			return nil
		},
	}
)

// 业务路由注册
func registerBizRouter(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.GET("/swagger", func(r *ghttp.Request) {
			r.Response.Write(consts.SwaggerUIPageContent)
		})

		group.Bind(
			controller.Address,
		)
	})
}
