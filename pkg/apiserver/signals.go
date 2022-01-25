/*
@Date: 2021/11/10 16:04
@Author: yvan.zhang
@File : signals
*/

package apiserver

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yvanz/gin-tmpl/pkg/logger"
)

const timeFormat = "0102150405"

var done = make(chan struct{})

func init() {
	go func() {
		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, os.Interrupt, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				dumpGoroutines()
			case syscall.SIGTERM, os.Interrupt:
				select {
				case <-done:
					// already closed
				default:
					close(done)
				}

				gracefulStop(signals)
			default:
				logger.Errorf("got unregistered signal: %+v", v)
			}
		}
	}()
}
