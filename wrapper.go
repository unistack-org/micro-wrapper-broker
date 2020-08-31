package broker

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/server"
)

func NewSubscribeWrapper(h broker.Handler, wrappers []server.SubscriberWrapper) broker.Handler {
	return func(evt broker.Event) error {
		fn := func(ctx context.Context, msg server.Message) error {
			md, ok := metadata.FromContext(ctx)
			if !ok {
				md = make(metadata.Metadata)
			}

			nmd := make(metadata.Metadata)
			if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
				opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.TextMap, opentracing.TextMapCarrier(nmd))
			} else if spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md)); err == nil {
				opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, opentracing.TextMapCarrier(nmd))
			}

			if len(nmd) > 0 {
				msg.(*event).AddHeader(nmd)
			}

			return h(msg.(*event).Event())
		}

		ctx := metadata.NewContext(context.Background(), evt.Message().Header)
		msg := &event{evt}
		for i := len(wrappers); i > 0; i-- {
			fn = wrappers[i-1](fn)
		}

		return fn(ctx, msg)
	}
}
