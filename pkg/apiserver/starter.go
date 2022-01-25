/*
@Date: 2021/11/10 17:01
@Author: yvan.zhang
@File : starter
*/

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yvanz/gin-tmpl/pkg/logger"
)

// StartOption defines the method to customize http.Server.
type StartOption func(srv *http.Server)

// StartHTTP starts a http server.
func StartHTTP(host string, port int, handler http.Handler, opts ...StartOption) error {
	listenAddr := fmt.Sprintf("%s:%d", host, port)

	return start(listenAddr, handler, func(srv *http.Server) error {
		return srv.ListenAndServe()
	}, opts...)
}

// StartHTTPS starts a https server.
func StartHTTPS(conf APIConfig, handler http.Handler, opts ...StartOption) error {
	listenAddr := fmt.Sprintf("%s:%d", conf.App.HostIP, conf.App.APIPort)

	return start(listenAddr, handler, func(srv *http.Server) error {
		return srv.ListenAndServeTLS(conf.App.CertFile, conf.App.KeyFile)
	}, opts...)
}

func start(listenAddr string, handler http.Handler, run func(*http.Server) error, opts ...StartOption) (err error) {
	server := &http.Server{
		Addr:    listenAddr,
		Handler: handler,
	}

	for _, opt := range opts {
		opt(server)
	}

	waitForCalled := AddWrapUpListener(func() {
		if e := server.Shutdown(context.Background()); e != nil {
			logger.Error(e)
		}
	})

	defer func() {
		if err != nil {
			if err == http.ErrServerClosed {
				waitForCalled()
			} else {
				logger.Errorf("http with an error: %s", err.Error())
			}
		}
	}()

	return run(server)
}
