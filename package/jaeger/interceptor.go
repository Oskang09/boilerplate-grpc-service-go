package jaeger

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"bitbucket.org/revenuemonster/monster-api/kit/grpc/errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

// ServiceLog :
func ServiceLog(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		method := info.FullMethod

		var span opentracing.Span

		spanContext, err := extractSpanContext(ctx, tracer)
		if err == nil {
			span = tracer.StartSpan(
				method,
				ext.RPCServerOption(spanContext),
			)
		} else {
			span = tracer.StartSpan(method)
		}

		defer span.Finish()

		rsp, errs := handler(opentracing.ContextWithSpan(injectSpanContext(ctx, span, tracer), span), req)

		requestData, err := json.Marshal(req)
		if err == nil {
			span.LogKV("grpc.request.body", cleanSensitiveData(string(requestData)))
		}

		responseData := ""

		if errs != nil {
			span.SetTag("error", true)
			responseData = errors.Parse(errs).String()
		} else {
			span.SetTag("error", false)
			responseJSON, err := json.Marshal(rsp)
			if err == nil {
				responseData = string(responseJSON)
			}
		}
		span.LogKV("grpc.response.body", cleanSensitiveData(string(responseData)))

		log.Printf(`{"time":"%s","endpoint":"%s"}`, time.Now().UTC().Format(time.RFC3339), method)

		return rsp, errs
	}
}
