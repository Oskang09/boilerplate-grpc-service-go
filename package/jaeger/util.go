package jaeger

import (
	"context"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tidwall/sjson"
	"google.golang.org/grpc/metadata"
)

// Metadata :
type Metadata struct {
	metadata.MD
}

// NewMetadata :
func NewMetadata(md metadata.MD) Metadata {
	return Metadata{md}
}

// Set :
func (w Metadata) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey :
func (w Metadata) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return tracer.Extract(opentracing.HTTPHeaders, NewMetadata(md))
}

func injectSpanContext(ctx context.Context, span opentracing.Span, tracer opentracing.Tracer) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	mdWriter := NewMetadata(md)
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, mdWriter)
	if err != nil {
		md = metadata.MD{}
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func cleanSensitiveData(jsonData string) string {
	jsonData, _ = sjson.Delete(jsonData, "pin")
	jsonData, _ = sjson.Delete(jsonData, "password")
	jsonData, _ = sjson.Delete(jsonData, "transData.cardNo")
	jsonData, _ = sjson.Delete(jsonData, "clientSecret")
	jsonData, _ = sjson.Delete(jsonData, "Data.transData.cardNo")

	return jsonData
}

// StartSpanFromContext :
func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, title string, startAt time.Time) opentracing.Span {
	var span opentracing.Span

	spanContext, err := extractSpanContext(ctx, tracer)
	if err == nil {
		span = tracer.StartSpan(
			title,
			ext.RPCServerOption(spanContext),
			opentracing.StartTime(startAt),
		)
	} else {
		span = tracer.StartSpan(title)
	}

	return span
}
