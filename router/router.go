package router

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Context struct {
	Message	 *message.Message
	Context  context.Context
	Topic 	 string
	MessageEngine *MessageEngine
	handlers HandlersChain
	index    int8
}

type HandlerFunc func(*Context) error

type HandlersChain []HandlerFunc

type RouterGroup struct {
	Handlers	HandlersChain
	handler		HandlersChain
	engine   	*MessageEngine
}

func (c *Context) Next() (err error) {
	c.index++
	for c.index < int8(len(c.handlers)) {
		err = c.handlers[c.index](c)
		c.index++
	}
	return err
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middleware...)
}


func (group *RouterGroup) AddRoutes(eventName string, handlers HandlerFunc) {
	group.handler = append(group.handler, handlers)
	group.AddRouteOne(eventName, group.handler...)
}

func (group *RouterGroup) AddRouteOne(eventName string, handlers ...HandlerFunc) {
	handlers = group.combineHandlers(handlers)

	fmt.Printf("%#v\n", handlers)

	group.engine.addRoute(eventName, handlers)
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}
