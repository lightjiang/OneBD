package OneBD

import (
	"github.com/lightjiang/OneBD/core"
	"github.com/lightjiang/OneBD/libs/handler"
	"github.com/lightjiang/OneBD/rfc"
	"github.com/lightjiang/OneBD/utils/log"
	"strings"
	"testing"
)

type testHandler struct {
	handler.BaseHandler
}

func (h *testHandler) Get() (interface{}, error) {
	h.Meta().SetHeader("a", "1")
	h.Meta().Write([]byte(h.Meta().RemoteAddr()))
	h.Meta().StreamWrite(strings.NewReader("asdasdasd"))
	return nil, nil
}

func TestNew(t *testing.T) {
	cfg := &core.Config{
		Host:           "0.0.0.0:4000",
		Charset:        "",
		TimeFormat:     "",
		PostMaxMemory:  0,
		TlsCfg:         nil,
		MaxConnections: 0,
	}
	cfg.LoggerLevel = log.DebugLevel
	cfg.BuildLogger()
	app := New(cfg)
	newH := func() core.Handler {
		app.Logger().Info("creating a handler")
		return &testHandler{}
	}
	app.Router().SubRouter("/asd/sss").Set("asd/xxx", newH, rfc.MethodGet, rfc.MethodPost)
	app.Router().SubRouter("/asd/sss").Set("asd/zzz", newH, rfc.MethodPost)
	app.Router().SubRouter("/sss/asd").Set("/123/int:username/str:username", newH)
	app.Router().SetNotFoundFunc(func(m core.Meta) {
		app.Logger().Info("checking 404 status")
		m.Write([]byte(m.RequestPath()))
	})
	app.Router().SetInternalErrorFunc(func(meta core.Meta) {
		app.Logger().Info("checking 500 status")
	})
	err := app.Run()
	t.Error(err)
}
