package whatsmeow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"whatsm/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"go.mau.fi/whatsmeow/store"
)

type Notify struct {
	host string
	path string
}

type DeviceInfo struct {
	ID           string `json:"id,omitempty"`
	LID          string `json:"lid,omitempty"`
	BusinessName string `json:"businessName"`
	PushName     string `json:"pushName"`
	Platform     string `json:"platform,omitempty"`
}

type Status int

const (
	StatusDefault               Status = iota // 0
	StatusConnected                           // 1
	StatusConnectFailed                       // 2
	StatusDisconnected                        // 3
	StatusLoggedIn                            // 4
	StatusLoggedOut                           // 5
	StatusPermanentDisconnected               // 6
	StatusInnerError                          // 7
	StatusPairSuccess                         // 8
	StatusPairFailed                          // 9
	StatusKeepAliveTimeout                    // 10
	StatusKeepAliveRestored                   // 11
	StatusHistorySync                         // 12
	StatusHistorySyncFinished                 // 13
	StatusOfflineSyncPreview                  // 14
	StatusOfflineSyncFinished                 // 15
	StatusTempBanned                          // 16
	StatusClientOutdated                      // 17
	StatusStreamReplaced                      // 18
)

func (s Status) String() string {
	var message string
	switch s {
	case StatusDefault:
		message = "ok"
	case StatusConnected:
		message = "connected"
	case StatusDisconnected:
		message = "disconnected"
	case StatusConnectFailed:
		message = "connect failed"
	case StatusLoggedIn:
		message = "client logged in"
	case StatusLoggedOut:
		message = "client logged out"
	case StatusPermanentDisconnected:
		message = "client permanent disconnected"
	case StatusInnerError:
		message = "inner error"
	case StatusPairSuccess:
		message = "pair success"
	case StatusPairFailed:
		message = "pair failed"
	case StatusKeepAliveTimeout:
		message = "client keep alive timeout"
	case StatusKeepAliveRestored:
		message = "client keep alive restored"
	case StatusHistorySync:
		message = "history sync"
	case StatusHistorySyncFinished:
		message = "history sync finished"
	case StatusOfflineSyncPreview:
		message = "offline sync preview"
	case StatusOfflineSyncFinished:
		message = "offline sync finished"
	case StatusTempBanned:
		message = "client temp banned"
	case StatusClientOutdated:
		message = "client outdated"
	case StatusStreamReplaced:
		message = "maybe logged in from another device, stream replaced"
	default:
		message = "unknown"
	}
	return message
}

type NotifyMessage struct {
	Phone      string     `json:"phone"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
	Event      Status     `json:"event"`
	Message    string     `json:"message,omitempty"`
	Timestamp  int64      `json:"timestamp"`
	TraceID    string     `json:"traceID,omitempty"`
}

func NewNotify(host, path string) *Notify {
	return &Notify{
		host: host,
		path: path,
	}
}

func (n *Notify) NotifyEvent(ctx context.Context, status Status, dev *store.Device, message string) {
	msg := NotifyMessage{
		Phone: dev.ID.User,
		DeviceInfo: DeviceInfo{
			ID:           dev.ID.ADString(),
			LID:          dev.LID.ADString(),
			BusinessName: dev.BusinessName,
			Platform:     dev.Platform,
			PushName:     dev.PushName,
		},
		Event:     status,
		Message:   fmt.Sprintf("%s: %s", status.String(), message),
		Timestamp: time.Now().UnixMilli(),
		TraceID:   guid.S(),
	}

	jsonStr, _ := json.Marshal(msg)
	g.Log(consts.LogicLog).Debugf(ctx, "notify event: %s", jsonStr)

	response, err := g.Client().ContentJson().Post(ctx, n.host+n.path, msg)
	if err != nil {
		g.Log(consts.LogicLog).Errorf(ctx, "failed to post notify message: %v", err)
		return
	}
	g.Log(consts.LogicLog).Debugf(ctx, "notify state response: %v", response.ReadAllString())
	defer response.Close()
}
