package whatsmeow

import (
	"fmt"
	"reflect"
	"whatsm/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type session struct {
	cli *whatsmeow.Client
	sw  *sWhats

	hooks []EventHook
}

// 事件处理函数
func (s *session) eventHandler(evt any) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "handle event type: %T", evt)

	switch v := evt.(type) {
	case *events.QR:
		s.handleQRCode(v)
	case *events.PairSuccess:
		s.handlePairSuccess(v)
	case *events.PairError:
		s.handlePairError(v)
	case *events.QRScannedWithoutMultidevice:
		s.handleQRScannedWithoutMultidevice(v)
	case *events.Connected:
		s.handleConnected(v)
	case *events.KeepAliveTimeout:
		s.handleKeepAliveTimeout(v)
	case *events.KeepAliveRestored:
		s.handleKeepAliveRestored(v)
	case events.PermanentDisconnect:
		s.handlePermanentDisconnect(v)
	case *events.ManualLoginReconnect:
		s.handleManualLoginReconnect(v)
	case *events.ConnectFailure:
		s.handleConnectFailure(v)
	// case *events.TemporaryBan:
	// 	s.handleTemporaryBan(v)
	case *events.LoggedOut:
		s.handleLoggedOut(v)
	//case *events.StreamReplaced:
	//	s.handleStreamReplaced(v)
	// case *events.ClientOutdated:
	// 	s.handleClientOutdated(v)
	//case *events.CATRefreshError:
	//	s.handleCATRefreshError(v)
	case *events.StreamError:
		s.handleStreamError(v)
	case *events.Disconnected:
		s.handleDisconnected(v)
	case *events.HistorySync:
		s.handleHistorySync(v)
	case *events.Receipt:
		s.handleReceipt(v)
	case *events.Message:
		s.handleMessage(v)
	case *events.ChatPresence:
		s.handleChatPresence(v)
	case *events.Presence:
		s.handlePresence(v)
	case *events.JoinedGroup:
		s.handleJoinedGroup(v)
	case *events.GroupInfo:
		s.handleGroupInfo(v)
	case *events.Picture:
		s.handlePicture(v)
	case *events.UserAbout:
		s.handleUserAbout(v)
	case *events.IdentityChange:
		s.handleIdentityChange(v)
	case *events.PrivacySettings:
		s.handlePrivacySettings(v)
	case *events.OfflineSyncPreview:
		s.handleOfflineSyncPreview(v)
	case *events.OfflineSyncCompleted:
		s.handleOfflineSyncCompleted(v)
	case *events.MediaRetry:
		s.handleMediaRetry(v)
	case *events.Blocklist:
		s.handleBlocklist(v)
	case *events.NewsletterJoin:
		s.handleNewsletterJoin(v)
	case *events.NewsletterLeave:
		s.handleNewsletterLeave(v)
	case *events.NewsletterMuteChange:
		s.handleNewsletterMuteChange(v)
	case *events.NewsletterLiveUpdate:
		s.handleNewsletterLiveUpdate(v)
		// appstate events
	case *events.Contact:
		s.handleContactEvent(v)
	case *events.PushName:
		s.handlePushNameEvent(v)
	case *events.BusinessName:
		s.handleBusinessNameEvent(v)
	case *events.Pin:
		s.handlePinEvent(v)
	case *events.Star:
		s.handleStarEvent(v)
	case *events.DeleteForMe:
		s.handleDeleteForMeEvent(v)
	case *events.Mute:
		s.handleMuteEvent(v)
	case *events.Archive:
		s.handleArchiveEvent(v)
	case *events.MarkChatAsRead:
		s.handleMarkChatAsReadEvent(v)
	case *events.ClearChat:
		s.handleClearChatEvent(v)
	case *events.DeleteChat:
		s.handleDeleteChatEvent(v)
	case *events.PushNameSetting:
		s.handlePushNameSettingEvent(v)
	case *events.UnarchiveChatsSetting:
		s.handleUnarchiveChatsSettingEvent(v)
	case *events.UserStatusMute:
		s.handleUserStatusMuteEvent(v)
	case *events.LabelEdit:
		s.handleLabelEditEvent(v)
	case *events.LabelAssociationChat:
		s.handleLabelAssociationChatEvent(v)
	case *events.LabelAssociationMessage:
		s.handleLabelAssociationMessageEvent(v)
	case *events.AppState:
		s.handleAppStateEvent(v)
	case *events.AppStateSyncComplete:
		s.handleAppStateSyncCompleteEvent(v)
	case *events.AppStateSyncError:
		s.handleAppStateSyncErrorEvent(v)
	default:
		actualType := reflect.TypeOf(v)
		g.Log(consts.LogicLog).Warningf(s.sw.ctx, "unHandle event: %s, user: %s", actualType, s.cli.Store.ID.ADString())
	}
}

// 处理二维码事件
func (s *session) handleQRCode(evt *events.QR) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: QR, codes: %v", evt.Codes)
}

// 处理配对成功事件
func (s *session) handlePairSuccess(evt *events.PairSuccess) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PairSuccess, deviceID: %s, businessName: %s, platform: %s", evt.ID, evt.BusinessName, evt.Platform)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusPairSuccess, s.cli.Store, "")
}

// 处理配对失败事件
func (s *session) handlePairError(evt *events.PairError) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PairError, error: %v", evt.Error)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusPairFailed, s.cli.Store, evt.Error.Error())
}

// 处理手机未启用多设备时的扫描二维码事件
func (s *session) handleQRScannedWithoutMultidevice(*events.QRScannedWithoutMultidevice) {
	g.Log(consts.LogicLog).Debug(s.sw.ctx, "event: QRScannedWithoutMultidevice, please enable multi-device and rescan")
}

// 处理客户端成功连接事件
func (s *session) handleConnected(evt *events.Connected) {
	// _ = service.Hook().Trigger(gctx.New(), &model.HookData{Event: consts.EventLogin, Phone: s.cli.Store.ID.User})
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Connected, user: %s", s.cli.Store.ID.ADString())
	s.sw.mu.Lock()
	s.sw.sessions[s.cli.Store.ID.User] = s
	s.sw.mu.Unlock()
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Connected, session added: %s", s.cli.Store.ID.ADString())
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusConnected, s.cli.Store, "")
}

// 处理心跳超时事件
func (s *session) handleKeepAliveTimeout(evt *events.KeepAliveTimeout) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: KeepAliveTimeout, errorCount: %d, lastSuccess: %v", evt.ErrorCount, evt.LastSuccess)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusKeepAliveTimeout, s.cli.Store, fmt.Sprintf("errorCount: %d, lastSuccess: %v", evt.ErrorCount, evt.LastSuccess))
}

// 处理心跳恢复事件
func (s *session) handleKeepAliveRestored(*events.KeepAliveRestored) {
	g.Log(consts.LogicLog).Debug(s.sw.ctx, "event: KeepAliveRestored, heartbeat restored")
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusKeepAliveRestored, s.cli.Store, "")
}

// 处理客户端永久断开事件
func (s *session) handlePermanentDisconnect(evt events.PermanentDisconnect) {

	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PermanentDisconnect, user: %s, description: %s", s.cli.Store.ID.ADString(), evt.PermanentDisconnectDescription())

	switch t := evt.(type) {
	case *events.LoggedOut:
		s.handleLoggedOut(t)
		return
	case *events.StreamReplaced:
		s.handleStreamReplaced(t)
		return
	case *events.ClientOutdated:
		s.handleClientOutdated(t)
		return
	case *events.CATRefreshError:
		s.handleCATRefreshError(t)
		return
	case *events.TemporaryBan:
		s.handleTemporaryBan(t)
		return
	default:
		s.sw.notify.NotifyEvent(s.sw.ctx, StatusPermanentDisconnected, s.cli.Store, evt.PermanentDisconnectDescription())
		return
	}
}

// 处理登出事件
func (s *session) handleLoggedOut(evt *events.LoggedOut) {
	// _ = service.Hook().Trigger(gctx.New(), &model.HookData{Event: consts.EventLogout, Phone: s.cli.Store.ID.User, Message: evt.PermanentDisconnectDescription()})
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: LoggedOut, user: %s, description: %s", s.cli.Store.ID.ADString(), evt.PermanentDisconnectDescription())
	s.sw.mu.Lock()
	delete(s.sw.sessions, s.cli.Store.ID.User)
	s.sw.mu.Unlock()
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: LoggedOut, session removed: %s", s.cli.Store.ID.ADString())
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusLoggedOut, s.cli.Store, evt.PermanentDisconnectDescription())
}

// 处理流被替换事件
func (s *session) handleStreamReplaced(evt *events.StreamReplaced) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: StreamReplaced, user: %s, description: %s", s.cli.Store.ID.ADString(), evt.PermanentDisconnectDescription())
	// s.sw.mu.Lock()
	// delete(s.sw.sessions, s.cli.Store.ID.User)
	// s.sw.mu.Unlock()
	// g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: StreamReplaced, session removed: %s", s.cli.Store.ID.ADString())
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusStreamReplaced, s.cli.Store, evt.PermanentDisconnectDescription())
}

// 处理手动登录重连事件
func (s *session) handleManualLoginReconnect(*events.ManualLoginReconnect) {
	g.Log(consts.LogicLog).Debug(s.sw.ctx, "event: ManualLoginReconnect, waiting for manual reconnect")
}

// 处理临时封禁事件
func (s *session) handleTemporaryBan(evt *events.TemporaryBan) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: TemporaryBan, reason: %v, expire: %v", evt.Code, evt.Expire)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusTempBanned, s.cli.Store, fmt.Sprintf("reason: %v, expire: %v", evt.Code, evt.Expire))
}

// 处理连接失败事件
func (s *session) handleConnectFailure(evt *events.ConnectFailure) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: ConnectFailure, error: %s", evt.Message)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusConnectFailed, s.cli.Store, fmt.Sprintf("reason: %s, error: %s", evt.Reason, evt.Message))
}

// 处理客户端过时事件
func (s *session) handleClientOutdated(evt *events.ClientOutdated) {
	g.Log(consts.LogicLog).Debug(s.sw.ctx, "event: ClientOutdated, client is outdated")
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusClientOutdated, s.cli.Store, evt.PermanentDisconnectDescription())
}

// 处理CAT刷新错误事件
func (s *session) handleCATRefreshError(evt *events.CATRefreshError) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: CATRefreshError, error: %v", evt.Error)
}

// 处理流错误事件
func (s *session) handleStreamError(evt *events.StreamError) {
	g.Log(consts.LogicLog).Warningf(s.sw.ctx, "event: StreamError, errorCode: %s", evt.Code)
}

// 处理客户端断开连接事件
func (s *session) handleDisconnected(*events.Disconnected) {
	// _ = service.Hook().Trigger(gctx.New(), &model.HookData{Event: consts.EventDisconnected, Phone: s.cli.Store.ID.User})
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Disconnected, user: %s", s.cli.Store.ID.ADString())
	s.sw.mu.Lock()
	delete(s.sw.sessions, s.cli.Store.ID.User)
	s.sw.mu.Unlock()
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Disconnected, session removed: %s", s.cli.Store.ID.ADString())
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusDisconnected, s.cli.Store, "")
}

// 处理历史同步事件
func (s *session) handleHistorySync(*events.HistorySync) {
	g.Log(consts.LogicLog).Debug(s.sw.ctx, "event: HistorySync, syncing historical messages")
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusHistorySync, s.cli.Store, "")
}

// 处理接收消息回执事件
func (s *session) handleReceipt(evt *events.Receipt) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Receipt, type: %s, messageIDs: %v", evt.Type, evt.MessageIDs)
}

// 处理接收消息事件
func (s *session) handleMessage(evt *events.Message) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Message, from %s, to: %s", evt.Info.Sender, evt.Message.GetConversation())
	for _, fn := range s.hooks {
		fn(s.sw.ctx, s.cli, evt)
	}
}

// 处理聊天状态事件
func (s *session) handleChatPresence(evt *events.ChatPresence) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: ChatPresence, sender: %s, state: %s", evt.Sender, evt.State)
}

// 处理用户在线状态事件
func (s *session) handlePresence(evt *events.Presence) {
	if evt.Unavailable {
		g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Presence, user: %s, status: offline", evt.From)
	} else {
		g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Presence, user: %s, lastSeen: %v", evt.From, evt.LastSeen)
	}
}

// 处理加入群组事件
func (s *session) handleJoinedGroup(evt *events.JoinedGroup) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: JoinedGroup, group: %s, reason: %s", evt.GroupInfo.Name, evt.Reason)
}

// 处理群组信息更新事件
func (s *session) handleGroupInfo(evt *events.GroupInfo) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: GroupInfo, group: %s, action: %v", evt.JID, evt.Name)
}

// 处理用户头像变更事件
func (s *session) handlePicture(evt *events.Picture) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Picture, user: %s, pictureID: %s", evt.JID, evt.PictureID)
}

// 处理用户状态更新事件
func (s *session) handleUserAbout(evt *events.UserAbout) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: UserAbout, user: %s, status: %s", evt.JID, evt.Status)
}

// 处理身份更改事件
func (s *session) handleIdentityChange(evt *events.IdentityChange) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: IdentityChange, user: %s, timestamp: %v", evt.JID, evt.Timestamp)
}

// 处理隐私设置更新事件
func (s *session) handlePrivacySettings(evt *events.PrivacySettings) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PrivacySettings, newSettings: %v", evt.NewSettings)
}

// 处理离线同步预览事件
func (s *session) handleOfflineSyncPreview(evt *events.OfflineSyncPreview) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: OfflineSyncPreview, total: %d", evt.Total)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusOfflineSyncPreview, s.cli.Store, fmt.Sprintf("total: %d", evt.Total))
}

// 处理离线同步完成事件
func (s *session) handleOfflineSyncCompleted(evt *events.OfflineSyncCompleted) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: OfflineSyncCompleted, count: %d", evt.Count)
	s.sw.notify.NotifyEvent(s.sw.ctx, StatusOfflineSyncFinished, s.cli.Store, fmt.Sprintf("count: %d", evt.Count))
}

// 处理媒体重试事件
func (s *session) handleMediaRetry(evt *events.MediaRetry) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: MediaRetry, messageID: %s", evt.MessageID)
}

// 处理黑名单事件
func (s *session) handleBlocklist(evt *events.Blocklist) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Blocklist, action: %s", evt.Action)
}

// 处理新闻订阅加入事件
func (s *session) handleNewsletterJoin(evt *events.NewsletterJoin) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: NewsletterJoin, metadata: %v", evt)
}

// 处理新闻订阅退出事件
func (s *session) handleNewsletterLeave(evt *events.NewsletterLeave) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: NewsletterLeave, role: %v", evt.Role)
}

// 处理新闻订阅静音变更事件
func (s *session) handleNewsletterMuteChange(evt *events.NewsletterMuteChange) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: NewsletterMuteChange, muteState: %v", evt.Mute)
}

// 处理新闻实时更新事件
func (s *session) handleNewsletterLiveUpdate(evt *events.NewsletterLiveUpdate) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: NewsletterLiveUpdate, messages: %v", evt.Messages)
}

// Contact event处理
func (s *session) handleContactEvent(v *events.Contact) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Contact, JID: %s, Timestamp: %v, Action: %v", v.JID, v.Timestamp, v.Action)
}

// PushName event处理
func (s *session) handlePushNameEvent(v *events.PushName) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PushName, JID: %s, OldPushName: %s, NewPushName: %s, Timestamp: %v", v.JID, v.OldPushName, v.NewPushName, v.Message.Timestamp)
}

// BusinessName event处理
func (s *session) handleBusinessNameEvent(v *events.BusinessName) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: BusinessName, JID: %s, OldBusinessName: %s, NewBusinessName: %s, Timestamp: %v", v.JID, v.OldBusinessName, v.NewBusinessName, v.Message.Timestamp)
}

// Pin event处理
func (s *session) handlePinEvent(v *events.Pin) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Pin, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// Star event处理
func (s *session) handleStarEvent(v *events.Star) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Star, ChatJID: %s, SenderJID: %s, MessageID: %s, Action: %v, Timestamp: %v", v.ChatJID, v.SenderJID, v.MessageID, v.Action, v.Timestamp)
}

// DeleteForMe event处理
func (s *session) handleDeleteForMeEvent(v *events.DeleteForMe) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: DeleteForMe, ChatJID: %s, SenderJID: %s, MessageID: %s, Action: %v, Timestamp: %v", v.ChatJID, v.SenderJID, v.MessageID, v.Action, v.Timestamp)
}

// Mute event处理
func (s *session) handleMuteEvent(v *events.Mute) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Mute, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// Archive event处理
func (s *session) handleArchiveEvent(v *events.Archive) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: Archive, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// MarkChatAsRead event处理
func (s *session) handleMarkChatAsReadEvent(v *events.MarkChatAsRead) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: MarkChatAsRead, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// ClearChat event处理
func (s *session) handleClearChatEvent(v *events.ClearChat) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: ClearChat, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// DeleteChat event处理
func (s *session) handleDeleteChatEvent(v *events.DeleteChat) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: DeleteChat, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// PushNameSetting event处理
func (s *session) handlePushNameSettingEvent(v *events.PushNameSetting) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: PushNameSetting, Action: %v, Timestamp: %v", v.Action, v.Timestamp)
}

// UnarchiveChatsSetting event处理
func (s *session) handleUnarchiveChatsSettingEvent(v *events.UnarchiveChatsSetting) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: UnarchiveChatsSetting, Action: %v, Timestamp: %v", v.Action, v.Timestamp)
}

// UserStatusMute event处理
func (s *session) handleUserStatusMuteEvent(v *events.UserStatusMute) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: UserStatusMute, JID: %s, Action: %v, Timestamp: %v", v.JID, v.Action, v.Timestamp)
}

// LabelEdit event处理
func (s *session) handleLabelEditEvent(v *events.LabelEdit) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: LabelEdit, LabelID: %s, Action: %v, Timestamp: %v", v.LabelID, v.Action, v.Timestamp)
}

// LabelAssociationChat event处理
func (s *session) handleLabelAssociationChatEvent(v *events.LabelAssociationChat) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: LabelAssociationChat, JID: %s, LabelID: %s, Action: %v, Timestamp: %v", v.JID, v.LabelID, v.Action, v.Timestamp)
}

// LabelAssociationMessage event处理
func (s *session) handleLabelAssociationMessageEvent(v *events.LabelAssociationMessage) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: LabelAssociationMessage, JID: %s, MessageID: %s, LabelID: %s, Action: %v, Timestamp: %v", v.JID, v.MessageID, v.LabelID, v.Action, v.Timestamp)
}

// AppState event处理
func (s *session) handleAppStateEvent(v *events.AppState) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: AppState, Index: %v, Action: %v", v.Index, v.SyncActionValue)
}

// AppStateSyncComplete event处理
func (s *session) handleAppStateSyncCompleteEvent(v *events.AppStateSyncComplete) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: AppStateSyncComplete, Name: %s, Version: %v, Recovery: %v", v.Name, v.Version, v.Recovery)
}

// AppStateSyncError event处理
func (s *session) handleAppStateSyncErrorEvent(v *events.AppStateSyncError) {
	g.Log(consts.LogicLog).Debugf(s.sw.ctx, "event: AppStateSyncError, Name: %s, Error: %v, FullSync: %v", v.Name, v.Error, v.FullSync)
}
