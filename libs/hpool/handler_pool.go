package hpool

import (
	"github.com/lightjiang/OneBD/core"
	"net/http"
	"sync"
)

type handlerPool struct {
	newFunc func() core.Handler
	pool    *sync.Pool
}

func New(newFunc func() core.Handler) core.HandlerPool {
	p := &handlerPool{
		pool: &sync.Pool{},
	}
	p.SetNew(newFunc)
	return p
}

func (p *handlerPool) SetNew(newFunc func() core.Handler) {
	p.newFunc = newFunc
	p.pool.New = func() interface{} {
		return newFunc()
	}
}

func (p *handlerPool) Acquire(w http.ResponseWriter, r *http.Request) core.Handler {
	ctx := p.pool.Get().(core.Handler)
	ctx.Init(w, r)
	return ctx
}

func (p *handlerPool) Release(h core.Handler) {
	h.TryReset()
	p.pool.Put(h)
}