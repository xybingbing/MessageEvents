package pub

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/xybingbing/MessageEvents/router"
	"strings"
)


func Tracer(ctx *router.Context) error {
	msg := ctx.Message
	EventName := msg.Metadata.Get("EventName")
	span, _ := opentracing.StartSpanFromContext(ctx.Context, "MessageEvents_Publish_" + ctx.Topic + "_" + EventName)
	if span != nil {
		span.LogKV("UUID", msg.UUID)
		span.LogKV("EventName", EventName)
		span.LogKV("Message", string(msg.Payload))
		defer span.Finish()
		//获取traceID并传递
		tracer := opentracing.GlobalTracer()
		carrier := opentracing.HTTPHeadersCarrier{}
		err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err == nil {
			for key, value := range carrier {
				if strings.ToLower(key) == jaeger.TraceContextHeaderName {
					msg.Metadata.Set(jaeger.TraceContextHeaderName, value[0])
					break;
				}
			}
		}
	}

	return ctx.Next()
}
