/*
@Date: 2021/11/10 10:55
@Author: yvanz
@File : server
*/

package apiserver

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yvanz/gin-tmpl/pkg/ginpprof"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"github.com/yvanz/gin-tmpl/pkg/middleware"
	"github.com/yvanz/gin-tmpl/pkg/tracer"
)

type Server struct {
	traceIO     io.Closer
	logger      *logger.DemoLog
	adminEngine *gin.Engine
	engine      *gin.Engine
	tracer      opentracing.Tracer
	conf        APIConfig
}

// CreateNewServer create a new server with gin
func CreateNewServer(ctx context.Context, c APIConfig, registerHandler func(opentracing.Tracer, *gin.Engine), opts ...ServerOption) *Server {
	server, err := newServer(ctx, c, registerHandler, opts)
	if err != nil {
		logger.Fatal(err)
	}

	return server
}

func newServer(ctx context.Context, c APIConfig, registerHandler func(opentracing.Tracer, *gin.Engine), options []ServerOption) (server *Server, err error) {
	opts := &serverOptions{}
	for _, o := range options {
		o(opts)
	}

	server = &Server{
		conf:   c,
		logger: c.buildLogger(),
	}

	// tracer 初始化必须在其他组件之前
	if c.Tracer.LocalAgentHostPort != "" {
		tra, cli, e := tracer.NewJaegerTracer(c.App.ServiceName, &c.Tracer, server.logger)
		if e != nil {
			err = e
			return
		}

		server.tracer = tra
		server.traceIO = cli
	}

	server.initGin(registerHandler)
	server.initAdmin()

	return server, c.initService(ctx, opts)
}

func (s *Server) initGin(registerHandler func(opentracing.Tracer, *gin.Engine)) {
	switch s.conf.App.RunMode {
	case RunModeRelease, RunModeProd, RunModeProduction:
		gin.SetMode(gin.ReleaseMode)
	case RunModeTest, RunModeDev:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	g := gin.New()
	g.Use(gin.Recovery(), middleware.GinFormatterLog(), middleware.Cors())

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ret_code": 0,
			"message":  "pong",
		})
	})

	if s.conf.App.RunMode == RunModeDebug {
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		gin.DisableConsoleColor()
	}

	if registerHandler != nil {
		registerHandler(s.GetTracer(), g)
	}

	s.engine = g
}

func (s *Server) initAdmin() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	g := gin.New()
	g.Use(middleware.GinFormatterLog(), gin.Recovery())

	ginpprof.Wrap(g)
	logger.Wrap(g)

	s.adminEngine = g
}

func (s *Server) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *Server) Start() {
	s.logger.Infof("starting server at %s: %d", s.conf.App.HostIP, s.conf.App.APIPort)

	go func() {
		s.logger.Infof("starting admin server at %s: %d", s.conf.App.HostIP, s.conf.App.AdminPort)
		err := StartHTTP(s.conf.App.HostIP, s.conf.App.AdminPort, s.adminEngine)
		handleError(err)
	}()

	var err error
	if s.conf.App.CertFile != "" && s.conf.App.KeyFile != "" {
		err = StartHTTPS(s.conf, s.engine)
	} else {
		err = StartHTTP(s.conf.App.HostIP, s.conf.App.APIPort, s.engine)
	}

	handleError(err)
}

func (s *Server) StartAdminOnly() {
	s.logger.Infof("starting admin server at %s: %d", s.conf.App.HostIP, s.conf.App.AdminPort)
	err := StartHTTP(s.conf.App.HostIP, s.conf.App.AdminPort, s.adminEngine)

	handleError(err)
}

func (s *Server) Stop() {
	_ = s.logger.Sync()

	if s.tracer != nil {
		_ = s.traceIO.Close()
	}
}

func handleError(err error) {
	// ErrServerClosed means the server is closed manually
	if err == nil || err == http.ErrServerClosed {
		return
	}

	logger.Fatal(err)
}
