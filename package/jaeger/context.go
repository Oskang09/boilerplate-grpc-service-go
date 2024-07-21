package jaeger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

// Context :
func Context(ctx context.Context, tracer opentracing.Tracer) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	sp := opentracing.SpanFromContext(ctx)

	mdWriter := NewMetadata(md)
	err := tracer.Inject(sp.Context(), opentracing.HTTPHeaders, mdWriter)
	if err != nil {
		md = metadata.MD{}
	}

	return metadata.NewOutgoingContext(ctx, md)
}
