/*
@Date: 2021/11/17 16:31
@Author: yvanz
@File : ctx
*/

package gadget

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

const SpanCtxKey = "span_ctx"

func ExtractTraceSpan(ctx context.Context) (spanCtx context.Context, err error) {
	if ctx == nil {
		return spanCtx, fmt.Errorf("ctx is nil")
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		if _, ok := span.Context().(jaeger.SpanContext); ok {
			return ctx, err
		}
	}

	getSpan := ctx.Value(SpanCtxKey)
	if getSpan == nil {
		return spanCtx, fmt.Errorf("span context not found")
	}

	spanCtx, ok := getSpan.(context.Context)
	if !ok {
		return spanCtx, fmt.Errorf("%s is not a context", SpanCtxKey)
	}

	return spanCtx, err
}
