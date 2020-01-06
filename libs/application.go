package libs

import (
	"crypto/tls"
	"github.com/lightjiang/OneBD/config"
	"github.com/lightjiang/OneBD/core"
	"go.uber.org/zap"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
)

type Application struct {
	config   *config.Config
	ctxPool  core.CtxPool
	server   *http.Server
	listener net.Listener
	router   core.Router
	// todo:: log 分级处理log
	// todo:: router
	//
}

func NewApplication(cfg *config.Config) *Application {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	app := &Application{
		config: cfg,
	}
	app.server = &http.Server{
		Addr:              cfg.Host,
		TLSConfig:         cfg.TlsCfg,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		// TODO
		BaseContext: nil,
		// TODO
		ConnContext: nil,
	}

	// 判断是否使用内置context
	newCtx := func() core.Context {
		return NewContext(app)
	}
	if cfg.NewCtx != nil {
		newCtx = cfg.NewCtx
	}

	// 判断是否使用内置ctx pool
	if cfg.CtxPool != nil {
		app.ctxPool = cfg.CtxPool
		app.ctxPool.SetCtx(newCtx)
	} else {
		app.ctxPool = NewCtxPool(newCtx)
	}

	// 判断是否使用内置路由
	if cfg.Router != nil {
		app.router = cfg.Router
	} else {
		app.router = NewRouter(app)
	}
	app.server.Handler = app.router
	return app
}

func (app *Application) Logger() *zap.Logger {
	return app.config.Logger
}

func (app *Application) Router() core.Router {
	return app.router
}

func (app *Application) Config() *config.Config {
	return app.config
}

func (app *Application) Run() error {
	l, e := app.netListener()
	if e != nil {
		return e
	}
	return app.server.Serve(l)
}

func (app *Application) netListener() (net.Listener, error) {
	if app.listener != nil {
		return app.listener, nil
	}
	l, err := net.Listen("tcp", app.config.Host)
	if err != nil {
		return nil, err
	}
	if app.config.TlsCfg != nil && len(app.config.TlsCfg.Certificates) > 0 && app.config.TlsCfg.GetCertificate != nil {
		l = tls.NewListener(l, app.config.TlsCfg)
	}
	if app.config.MaxConnections > 0 {
		l = netutil.LimitListener(l, app.config.MaxConnections)
	}
	app.listener = l
	return app.listener, nil
}