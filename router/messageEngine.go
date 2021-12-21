package router

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"sync"
)

type ErrFail = error

var (
	ErrRetryFail = fmt.Errorf("重试3次失败")
)

type MessageEngine struct{
	RouterGroup
	tree	map[string]HandlersChain
	pool	sync.Pool
}

func NewMessageEngine() *MessageEngine {
	engine := &MessageEngine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			handler: nil,
		},
		tree: make(map[string]HandlersChain),
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

func (engine *MessageEngine) allocateContext() *Context {
	return &Context{}
}

//获取路由下的相关HandlersChain
func (engine *MessageEngine) getValue(eventName string)(handlers HandlersChain){
	handlers, ok := engine.tree[eventName]
	if !ok {
		return nil
	}
	return
}

func (engine *MessageEngine) addRoute(eventName string, handlers HandlersChain) {
	engine.tree[eventName] = handlers
}

func (engine *MessageEngine) Use(middleware ...HandlerFunc)  {
	engine.RouterGroup.Use(middleware...)
}

func (engine *MessageEngine) handle(c *Context) {
	eventName := c.Message.Metadata.Get("EventName")
	handlers := engine.getValue(eventName)
	if handlers != nil {
		c.handlers = handlers
		_ = c.Next()
		return
	}
}

func (engine *MessageEngine) Run(ctx context.Context, msg *message.Message, topic string){
	c := engine.pool.Get().(*Context)
	c.Context = ctx
	c.Message = msg
	c.Topic = topic
	c.MessageEngine = engine
	c.reset()
	engine.handle(c)
	engine.pool.Put(c)
}
