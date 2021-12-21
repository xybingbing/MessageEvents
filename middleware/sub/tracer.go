package sub

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/xybingbing/MessageEvents/router"
)

func Tracer(ctx *router.Context) error {
	msg := ctx.Message
	EventName := msg.Metadata.Get("EventName")

	//链路追踪
	tracer := opentracing.GlobalTracer()
	carrier := opentracing.HTTPHeadersCarrier{}
	carrier.Set(jaeger.TraceContextHeaderName, msg.Metadata.Get(jaeger.TraceContextHeaderName))
	wireContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		serverSpan := opentracing.StartSpan("MessageEvents_Subscriber_" + ctx.Topic + EventName, ext.RPCServerOption(wireContext))
		opentracing.ContextWithSpan(ctx.Context, serverSpan)
		if serverSpan != nil {
			serverSpan.LogKV("UUID", msg.UUID)
			serverSpan.LogKV("EventName", EventName)
			serverSpan.LogKV("Message", string(msg.Payload))
			serverSpan.Finish()
		}
	}

	return ctx.Next()
}