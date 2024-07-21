package grpc

import (
	"fmt"
	"net"
	"os"

	jg "bitbucket.org/revenuemonster/monster-api/kit/jaeger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server :
type Server struct {
	*g.Server
}

// NewServer :
func NewServer(isNoobTracer bool) *Server {
	initOpenTracingLogging(isNoobTracer)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(jg.ServiceLog(opentracing.GlobalTracer())),
	)
	reflection.Register(grpcServer)

	return &Server{grpcServer}
}

// Start :
func (s *Server) Start(portNo string) error {
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNo)
	if err != nil {
		panic(err)
	}

	os.Stdout.WriteString("Starting [service]" + " : " + portNo + "\n")

	return s.Serve(lis)
}

func initOpenTracingLogging(isNoobTracer bool) {
	var tracer opentracing.Tracer

	if isNoobTracer {
		tracer = opentracing.NoopTracer{}
	} else {
		def := config.Configuration{
			Disabled: false,
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: true,
			},
		}

		cfg, err := def.FromEnv()
		if err != nil {
			panic(err)
		}

		jLogger := jaeger.StdLogger
		jMetricsFactory := metrics.NullFactory

		tracer, _, err = cfg.NewTracer(
			config.Logger(jLogger),
			config.Metrics(jMetricsFactory),
		)
		if err != nil {
			panic(fmt.Sprintf("error: cannot init Jaeger: %v\n", err))
		}
	}

	opentracing.SetGlobalTracer(tracer)
}
