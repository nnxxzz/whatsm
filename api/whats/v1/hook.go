package v1

import "github.com/gogf/gf/v2/frame/g"

type HookReq struct {
	g.Meta  `path:"/whats/hook" tags:"whats" method:"post" sm:"hook测试" dc:"hook测试"`
	Event   string `json:"event" v:"required" dc:"事件，1 Login, 2 Logout, 3 Disconnected 4 EventPermanentDisconnected"`
	Phone   string `json:"phone" dc:"账号"`
	Message string `json:"message" dc:"事件信息"`
}

type HookRes struct{}
