// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"whatsm/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IWhats interface {
		GetContactInfo(ctx context.Context, u string, t string) error
		Upload(ctx context.Context, from string, typ int, uf *ghttp.UploadFile) (*model.UploadOutput, error)
		SendTextMessage(ctx context.Context, from string, to string, msg string) error
		SendGroupTextMessage(ctx context.Context, from string, to string, msg string) error
		SendMediaMessage(ctx context.Context, in *model.SendMediaMessageInput) error
		// 检查账号是否登录
		IsWhatsAccountLogin(ctx context.Context, phone string) bool
		// 获取所有已登录的账号
		LoggedInAccounts() []string
		// Init connect to db
		Init(ctx context.Context) error
		Logout(ctx context.Context, phone string) error
		// create new device&session
		LoginPair(ctx context.Context, in *model.LoginPairInput) (*model.LoginPairOutput, error)
		RecoverSessions()
	}
)

var (
	localWhats IWhats
)

func Whats() IWhats {
	if localWhats == nil {
		panic("implement not found for interface IWhats, forgot register?")
	}
	return localWhats
}

func RegisterWhats(i IWhats) {
	localWhats = i
}
