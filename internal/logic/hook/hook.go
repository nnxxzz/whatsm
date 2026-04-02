package hook

import (
	"context"
	"whatsm/internal/consts"
	"whatsm/internal/model"
	"whatsm/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type sHook struct {
	c *gclient.Client
}

func init() {
	service.RegisterHook(New())
}

func New() service.IHook {
	return &sHook{
		c: gclient.New().ContentJson(),
	}
}

func (h *sHook) Trigger(ctx context.Context, data *model.HookData) error {
	gv, err := g.Cfg().Get(ctx, "callback.urls")
	if err != nil {
		return gerror.Wrap(err, "get callback.urls failed")
	}
	urls := gv.Strings()
	for _, url := range urls {
		g.Log(consts.LogicLog).Debugf(ctx, "trigger hook event: %d, phone: %s, message: %s, callback url: %s", data.Event, data.Phone, data.Message, url)
		if _, err := h.c.ContentJson().Post(ctx, url, data); err != nil {
			return gerror.Wrapf(err, "call back url: %s failed", url)
		}
	}
	//h.c.Post(ctx)
	return nil
}
