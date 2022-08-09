package ginpprof

import (
	"expvar"
	"fmt"
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Wrap(router *gin.Engine) {
	WrapGroup(&router.RouterGroup)
}

func WrapGroup(router *gin.RouterGroup) {
	routers := []struct {
		Handler gin.HandlerFunc
		Method  string
		Path    string
	}{
		{IndexHandler(), "GET", "/debug/pprof/"},
		{HeapHandler(), "GET", "/debug/pprof/heap"},
		{GoroutineHandler(), "GET", "/debug/pprof/goroutine"},
		{BlockHandler(), "GET", "/debug/pprof/block"},
		{ThreadCreateHandler(), "GET", "/debug/pprof/threadcreate"},
		{CmdlineHandler(), "GET", "/debug/pprof/cmdline"},
		{ProfileHandler(), "GET", "/debug/pprof/profile"},
		{SymbolHandler(), "GET", "/debug/pprof/symbol"},
		{SymbolHandler(), "POST", "/debug/pprof/symbol"},
		{TraceHandler(), "GET", "/debug/pprof/trace"},
		{MutexHandler(), "GET", "/debug/pprof/mutex"},
		{ExpvarHandler(), "GET", "/debug/vars"},
		{PromhttpHandler(), "GET", "/metrics"},
	}

	basePath := strings.TrimSuffix(router.BasePath(), "/")
	var prefix string

	switch {
	case basePath == "":
		prefix = ""
	case strings.HasSuffix(basePath, "/debug"):
		prefix = "/debug"
	case strings.HasSuffix(basePath, "/debug/pprof"):
		prefix = "/debug/pprof"
	}

	prometheus.EnableHandlingTimeHistogram()
	prometheus.EnableClientHandlingTimeHistogram()

	for _, r := range routers {
		router.Handle(r.Method, strings.TrimPrefix(r.Path, prefix), r.Handler)
	}
}

func IndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Index(c.Writer, c.Request)
	}
}

func HeapHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
	}
}

func GoroutineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
	}
}

func BlockHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
	}
}

func ThreadCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
	}
}

func CmdlineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// pprof.Handler("cmdline").ServeHTTP(c.Writer, c.Request)
		pprof.Cmdline(c.Writer, c.Request)
	}
}

func ProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// pprof.Handler("profile").ServeHTTP(c.Writer, c.Request)
		pprof.Profile(c.Writer, c.Request)
	}
}

func SymbolHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// pprof.Handler("symbol").ServeHTTP(c.Writer, c.Request)
		pprof.Symbol(c.Writer, c.Request)
	}
}

func TraceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// pprof.Handler("trace").ServeHTTP(c.Writer, c.Request)
		pprof.Trace(c.Writer, c.Request)
	}
}

func MutexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
	}
}

func PromhttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

func ExpvarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(c.Writer, "{\n")
		first := true
		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				fmt.Fprintf(c.Writer, ",\n")
			}
			first = false
			fmt.Fprintf(c.Writer, "%q: %s", kv.Key, kv.Value)
		})
		fmt.Fprintf(c.Writer, "\n}\n")
	}
}
