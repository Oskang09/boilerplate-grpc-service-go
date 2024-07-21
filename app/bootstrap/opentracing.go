package bootstrap

import (
	"fmt"
	"service/app/config"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	c "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func (bs *Bootstrap) initOpenTracing() *Bootstrap {
	var tracer opentracing.Tracer

	if config.IsLocal() {
		tracer = opentracing.NoopTracer{}
	} else {
		def := c.Configuration{
			ServiceName: "{service-name}",
			Disabled:    false,
			Sampler: &c.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &c.ReporterConfig{
				LogSpans:            true,
				BufferFlushInterval: 1 * time.Second,
			},
		}

		cfg, err := def.FromEnv()
		if err != nil {
			panic("Could not parse Jaeger env vars: " + err.Error())
		}

		jLogger := log.StdLogger
		jMetricsFactory := metrics.NullFactory

		tracer, _, err = cfg.NewTracer(
			c.Logger(jLogger),
			c.Metrics(jMetricsFactory),
			c.MaxTagValueLength(2048),
		)
		if err != nil {
			panic(fmt.Sprintf("Cannot init Jaeger: %v\n", err))
		}
	}

	opentracing.SetGlobalTracer(tracer)
	return bs
}
