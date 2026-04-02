package cmd

import (
	"context"
	"time"
	"whatsm/internal/consts"
	"whatsm/internal/controller/whats"
	"whatsm/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gconv"
)

type Main struct {
	g.Meta `name:"main" brief:"start server"`
}

type cStartInput struct {
	g.Meta  `name:"start" brief:"start whatsapp web server"`
	CfgFile string `short:"c" name:"config" arg:"true" brief:"config file"`
}

type CStartOutput struct{}

// func (m *Main) WaitToRecover(ctx context.Context, s *ghttp.Server) {

// }

func (m *Main) Start(ctx context.Context, in cStartInput) (out *CStartOutput, err error) {
	if err = service.Whats().Init(ctx); err != nil {
		return nil, err
	}
	if in.CfgFile != "" {
		g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(in.CfgFile)
	}

	g.Log(consts.LogicLog).Debugf(ctx, "starting server after init")

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Middleware(ghttp.MiddlewareCORS)
		group.Bind(
			whats.NewV1(),
		)
	})
	e, _ := g.Cfg().Get(ctx, "swagger.enabled")
	if gconv.Bool(e) {
		enhanceOpenAPIDoc(s)
	}

	// m.WaitToRecover(ctx, s)

	ready := make(chan struct{})

	go func() {
		for {
			port := s.GetListenedPort()
			if port > 0 {
				close(ready)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		<-ready
		g.Log(consts.LogicLog).Info(ctx, "server started, now to recover all sessions")
		service.Whats().RecoverSessions()
	}()

	s.Run()
	return nil, nil
}

func enhanceOpenAPIDoc(s *ghttp.Server) {
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`

	// API description.
	openapi.Info = goai.Info{
		Title:       consts.OpenApiTitle,
		Description: consts.OpenApiDesc,
	}
	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")
}
