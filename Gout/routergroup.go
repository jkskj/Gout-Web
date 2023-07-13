package Gout

import (
	"log"
	"net/http"
)

type RouterGroup struct {
	basePath string
	handlers HandlersChain // 中间件组
	engine   *Engine
}

// Group 路由分组
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	newGroup := &RouterGroup{
		basePath: group.basePath + relativePath,
		engine:   group.engine,
		handlers: group.combineHandlers(handlers),
	}
	return newGroup
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handlers HandlersChain) {
	pattern := group.basePath + comp
	log.Printf("Route %4s - %s", method, pattern)
	handlers = group.combineHandlers(handlers)
	group.engine.router.addRoute(method, pattern, handlers)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodGet, pattern, handlers)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodPost, pattern, handlers)
}

// PUT defines the method to add Put request
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodPut, relativePath, handlers)
}

// DELETE defines the method to add Delete request
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodDelete, relativePath, handlers)
}

// Use 使用中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.handlers = append(group.handlers, middlewares...)
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)
	return mergedHandlers
}
