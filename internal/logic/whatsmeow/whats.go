package whatsmeow

import (
	"context"
	"errors"
	"sync"
	"whatsm/internal/consts"
	"whatsm/internal/model"
	"whatsm/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type sWhats struct {
	c   *sqlstore.Container
	l   *logger
	ctx context.Context

	mu       sync.Mutex
	sessions map[string]*session
	notify   *Notify
}

func init() {
	service.RegisterWhats(New())
}

func New() service.IWhats {
	return &sWhats{sessions: make(map[string]*session), mu: sync.Mutex{}}
}

// 检查账号是否登录
func (s *sWhats) IsWhatsAccountLogin(ctx context.Context, phone string) bool {
	sess, ok := s.sessions[phone]
	if !ok {
		return false
	}
	return sess.cli.IsLoggedIn()
}

// 获取所有已登录的账号
func (s *sWhats) LoggedInAccounts() []string {
	phones := make([]string, 0)
	for p, sess := range s.sessions {
		if !sess.cli.IsLoggedIn() {
			continue
		}
		phones = append(phones, p)
	}
	return phones
}

// Init connect to db
func (s *sWhats) Init(ctx context.Context) error {
	s.ctx = ctx
	s.l = &logger{ctx: ctx}
	dialect := consts.DbDialectDefault
	address := consts.DbAddressDefault

	parser, _ := gcmd.Parse(g.MapStrBool{
		"c,config": true,
	})

	configPath := parser.GetOpt("c").String()
	if configPath != "" {
		g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(configPath)
		g.Log(consts.LogicLog).Infof(ctx, "config file set to: %s", configPath)
	}

	if dialectCfg, err := g.Cfg().Get(ctx, consts.DbDialectConfigKey); err == nil {
		dialect = dialectCfg.String()
	}
	if addressCfg, err := g.Cfg().Get(ctx, consts.DbAddressConfigKey); err == nil {
		address = addressCfg.String()
	}
	container, err := sqlstore.New(ctx, dialect, address, s.l)
	if err != nil {
		return gerror.Wrapf(err, "connect to db failed")
	}
	s.c = container

	host := consts.NotifyHostDefault
	path := consts.NotifyPathDefault
	if hostCfg, err := g.Cfg().Get(ctx, "callback.host"); err == nil {
		host = hostCfg.String()
	}
	if pathCfg, err := g.Cfg().Get(ctx, "callback.path"); err == nil {
		path = pathCfg.String()
	}
	g.Log(consts.LogicLog).Infof(ctx, "notify callback set to: %s%s", host, path)
	s.notify = NewNotify(host, path)
	return nil
}

func (s *sWhats) RecoverSessions() {
	g.Log(consts.LogicLog).Info(s.ctx, "recover sessions start")
	allDevices, err := s.c.GetAllDevices(s.ctx)
	if err != nil {
		g.Log(consts.LogicLog).Errorf(s.ctx, "get all devices failed, err: %v", err)
		return
	}

	if len(allDevices) == 0 {
		g.Log(consts.LogicLog).Info(s.ctx, "no device found, skip recover")
		return
	}

	devices := make(map[string]*store.Device, len(allDevices))
	for _, device := range allDevices {
		if dev, ok := devices[device.ID.User]; ok {
			if dev.ID.Device < device.ID.Device {
				devices[device.ID.User] = device
			}
		} else {
			devices[device.ID.User] = device
		}
	}

	for _, device := range devices {
		if err := s.autoLogin(s.ctx, device); err != nil {
			g.Log(consts.LogicLog).Errorf(s.ctx, "auto login failed for device %s, err: %v", device.ID.ADString(), err)
		} else {
			g.Log(consts.LogicLog).Infof(s.ctx, "auto login success for device %s", device.ID.ADString())
		}
	}

	g.Log(consts.LogicLog).Info(s.ctx, "recover sessions done")
}

func (s *sWhats) Logout(ctx context.Context, phone string) error {
	g.Log(consts.LogicLog).Debugf(ctx, "user %s logout", phone)
	if _, ok := s.sessions[phone]; !ok {
		return nil
	}
	return s.sessions[phone].cli.Logout(ctx)
}

func (s *sWhats) getDevice(ctx context.Context, phone string) (*store.Device, bool) {
	allDevices, err := s.c.GetAllDevices(ctx)
	if err != nil {
		g.Log(consts.LogicLog).Errorf(ctx, "get all devices failed, err: %v", err)
		return nil, false
	}

	for i := len(allDevices) - 1; i >= 0; i-- {
		device := allDevices[i]
		if device.ID.User == phone {
			g.Log(consts.LogicLog).Debugf(ctx, "get device, jid: %s, lid: %s, businessName: %s, pushName: %s", device.ID.ADString(), device.LID.ADString(), device.BusinessName, device.PushName)
			return device, true
		}
	}
	return nil, false
}

// create new device&session
func (s *sWhats) LoginPair(ctx context.Context, in *model.LoginPairInput) (*model.LoginPairOutput, error) {
	g.Log(consts.LogicLog).Debugf(ctx, "client loginPair, phone: %s, proxy: %s", in.Phone, in.Proxy)
	limit := consts.MaxUserDefault
	if maxUser, err := g.Cfg().Get(ctx, consts.MaxUserConfigKey); err == nil {
		if maxUser.Int() != 0 {
			limit = maxUser.Int()
		}
	}
	if len(s.sessions) >= limit {
		return nil, gerror.NewCode(gcode.New(1001, "login users limit", nil))
	}
	// jid := types.NewADJID(in.Phone, 0, 9)
	// g.Log(consts.LogicLog).Debugf(ctx, "login jid: %s", jid.ADString())
	// st, err := s.c.GetDevice(ctx, jid)
	// if err != nil {
	// 	g.Log(consts.LogicLog).Errorf(ctx, "get device failed, err: %v", err)
	// 	return nil, gerror.Wrapf(err, "device not found")
	// }

	st, ok := s.getDevice(ctx, in.Phone)
	if !ok {
		g.Log(consts.LogicLog).Debugf(ctx, "device not found, create new device")
		st = s.c.NewDeviceWithProxy(in.Proxy)
		// g.Log(consts.LogicLog).Debugf(ctx, "new device, jid: %s, lid: %s, businessName: %s, pushName: %s", st.ID.ADString(), st.LID.ADString(), st.BusinessName, consts.PushNameDefault)
	}
	client := whatsmeow.NewClient(st, s.l)
	if in.Proxy != "" {
		if err := client.SetProxyAddress(in.Proxy); err != nil {
			return nil, gerror.Wrapf(err, "set proxy address failed")
		}
	}
	sess := &session{cli: client, sw: s}
	autoMarkMessage := false
	if mark, err := g.Cfg().Get(ctx, consts.AutoMarkMessageKey); err == nil {
		autoMarkMessage = mark.Bool()
	}
	if autoMarkMessage {
		sess.hooks = []EventHook{HookMarkMessageAdRead}
	}
	client.AddEventHandler(sess.eventHandler)

	if client.Store.ID != nil {
		if err := client.Connect(); err != nil {
			g.Log(consts.LogicLog).Errorf(ctx, "client connect to whats server failed, err: %v", err)
			return nil, gerror.Wrapf(err, "client connect to whats server failed")
		}
		return &model.LoginPairOutput{}, nil
	}
	client.Store.Platform = consts.PlatformDefault
	if pfCfg, err := g.Cfg().Get(ctx, consts.PlatformConfigKey); err == nil {
		client.Store.Platform = pfCfg.String()
	}
	client.Store.BusinessName = consts.BusinessNameDefault
	if bnCfg, err := g.Cfg().Get(ctx, consts.BusinessNameConfigKey); err == nil {
		client.Store.BusinessName = bnCfg.String()
	}
	client.Store.PushName = consts.PushNameDefault
	if pnCfg, err := g.Cfg().Get(ctx, consts.PushNameConfigKey); err == nil {
		client.Store.PushName = pnCfg.String()
	}
	qrChan, _ := client.GetQRChannel(context.Background())
	if err := client.Connect(); err != nil {
		return nil, gerror.Wrapf(err, "client dial to whats server failed")
	}
	// ensure websocket is ok
	qrCode := <-qrChan

	code, err := client.PairPhone(ctx, in.Phone, true, whatsmeow.PairClientChrome, consts.ClientDisplayNameDefault)
	if err != nil {
		if errors.Is(err, whatsmeow.ErrIQRateOverLimit) {
			return nil, gerror.NewCode(gcode.New(1000, err.Error(), nil))
		}
		return nil, gerror.Wrapf(err, "create pair code failed")
	}
	// 等待连接成功
	return &model.LoginPairOutput{Code: code, QrCode: qrCode.Code}, nil
}

func (s *sWhats) autoLogin(ctx context.Context, device *store.Device) error {
	if device == nil {
		g.Log(consts.LogicLog).Errorf(ctx, "auto login device error, device is null")
		return gerror.Newf("device is null")
	}

	g.Log(consts.LogicLog).Debugf(ctx, "auto login device, jid: %s, lid: %s, businessName: %s, pushName: %s", device.ID.ADString(), device.LID.ADString(), device.BusinessName, device.PushName)

	client := whatsmeow.NewClient(device, s.l)
	if device.Proxy != "" {
		if err := client.SetProxyAddress(device.Proxy); err != nil {
			return gerror.Wrapf(err, "set proxy address failed")
		}
	}
	sess := &session{cli: client, sw: s}
	autoMarkMessage := false
	if mark, err := g.Cfg().Get(ctx, consts.AutoMarkMessageKey); err == nil {
		autoMarkMessage = mark.Bool()
	}
	if autoMarkMessage {
		sess.hooks = []EventHook{HookMarkMessageAdRead}
	}
	client.AddEventHandler(sess.eventHandler)

	if err := client.Connect(); err != nil {
		g.Log(consts.LogicLog).Errorf(ctx, "client connect to whats server failed, err: %v", err)
		return gerror.Wrapf(err, "client connect to whats server failed")
	}
	return nil

}
