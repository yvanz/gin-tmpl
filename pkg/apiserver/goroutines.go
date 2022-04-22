/*
@Date: 2021/11/10 16:10
@Author: yvanz
@File : goroutines
*/

package apiserver

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/yvanz/gin-tmpl/pkg/logger"
)

const (
	goroutineProfile = "goroutine"
	debugLevel       = 2
)

func dumpGoroutines() {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	dumpFile := path.Join(os.TempDir(), fmt.Sprintf("%s-%d-goroutines-%s.dump",
		command, pid, time.Now().Format(timeFormat)))

	logger.Infof("got dump goroutine signal, printing goroutine profile to %s", dumpFile)

	if f, err := os.Create(dumpFile); err != nil {
		logger.Errorf("failed to dump goroutine profile, error: %s", err.Error())
	} else {
		defer f.Close()
		_ = pprof.Lookup(goroutineProfile).WriteTo(f, debugLevel)
	}
}
