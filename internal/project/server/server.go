/*
@Date: 2021/1/12 下午2:27
@Author: yvan.zhang
@File : server
@Desc:
*/

package server

import (
	"context"
	"fmt"
	"gin-tmpl/api"
	"gin-tmpl/internal/config"
	"gin-tmpl/internal/queue"
	"gin-tmpl/models"
	"gin-tmpl/pkg/ginpprof"
	"gin-tmpl/pkg/kafka"
	"gin-tmpl/pkg/logger"
	"gin-tmpl/pkg/middleware"
	"gin-tmpl/pkg/xormmysql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Server struct {
	host              string
	conf              *config.Config
	ctx               context.Context
	log               *logger.DemoLog
	apiServer         *http.Server
	certFile, keyFile string
}

func New(ctx context.Context, conf *config.Config, log *logger.DemoLog) *Server {
	return &Server{
		host: conf.App.LocalIP,
		conf: conf,
		ctx:  ctx,
		log:  log,
	}
}

func (s *Server) Init() error {
	server := new(http.Server)
	server.Addr = fmt.Sprintf("%s:%d", s.host, s.conf.App.APIPort)
	server.Handler = s.handler()
	s.apiServer = server

	if err := s.initTable(); err != nil {
		return err
	}

	kafkaConsumer, err := kafka.Default().NewConsumer()
	if err != nil {
		return err
	}
	if err := queue.RunConsume(kafkaConsumer); err != nil {
		return err
	}

	kafkaProducer, err := kafka.Default().NewAsyncProducerClient()
	if err != nil {
		return err
	}
	queue.NewProducer(&kafkaProducer)

	return nil
}

func (s *Server) handler() http.Handler {
	g := gin.New()

	g.Use(middleware.GinFormatterLog())
	g.Use(gin.Recovery())
	g.Use(middleware.Cors())
	ginpprof.Wrap(g)
	v1Group := g.Group("/api")

	if s.conf.App.RunMode == `debug` {
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		gin.DisableConsoleColor()
	}

	api.Router(v1Group, s.log)
	return g
}

func (s *Server) initTable() error {
	return xormmysql.My().Master().Sync2(
		models.TableDemo,
	)
}

func (s *Server) Start() {
	logger.Infof("starting server: %s", s.apiServer.Addr)

	go func() {
		if s.certFile != "" && s.keyFile != "" {
			if err := s.apiServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && err != http.ErrServerClosed {
				s.log.Errorf("%s", err.Error())
				os.Exit(1)
			}
		} else {
			if err := s.apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.log.Errorf("%s", err.Error())
				os.Exit(1)
			}
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	if err := s.apiServer.Shutdown(ctx); err != nil {
		s.log.Errorf("failed to stop server, error: %s", err.Error())
	}
}
