/*
@Date: 2021/10/25 11:54
@Author: yvanz
@File : tracer
*/

package tracer

import (
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

var (
	// 仅用于单例模式
	_default opentracing.Tracer
	_dCloser io.Closer
)

type Config struct {
	BufferFlushInterval int    `yaml:"buffer_flush_interval"`
	LocalAgentHostPort  string `yaml:"local_agent_host_port" env:"TraceAgent" env-description:"host and port of jaeger agent"`
	LogSpan             bool   `yaml:"log_span" env:"TraceLog" env-description:"enable record span or not"`
}

func NewJaegerTracer(serviceName string, c *Config, logg *logger.DemoLog) (tra opentracing.Tracer, closer io.Closer, err error) {
	if c.LocalAgentHostPort == "" {
		return tra, closer, fmt.Errorf("no local agent host specified")
	}
	if _default != nil {
		return
	}

	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            c.LogSpan,
			BufferFlushInterval: time.Duration(c.BufferFlushInterval) * time.Second,
			LocalAgentHostPort:  c.LocalAgentHostPort,
		},
		Headers: &jaeger.HeadersConfig{
			JaegerDebugHeader:        "cmp-debug-id",
			JaegerBaggageHeader:      "cmp-baggage",
			TraceContextHeaderName:   "cmp-trace-id",
			TraceBaggageHeaderPrefix: "cmp-ctx",
		},
	}

	_default, _dCloser, err = cfg.NewTracer(config.Logger(logg))
	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(_default)
	return _default, _dCloser, err
}

func Default() opentracing.Tracer {
	return _default
}

func Span(serviceName string) opentracing.Span {
	span := Default().StartSpan(serviceName)
	return span
}
