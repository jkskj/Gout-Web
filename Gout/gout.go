package Gout

import (
	"net/http"
	"sync"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	pool   sync.Pool
}

type HandlersChain []HandlerFunc

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

func (engine *Engine) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	return http.ListenAndServe(address, engine)
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.reset(w, req)
	engine.router.HTTPRequest(c)
	engine.pool.Put(c)
}
