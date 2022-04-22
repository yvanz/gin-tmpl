/*
@Date: 2021/11/10 16:12
@Author: yvanz
@File : shutdown
*/

package apiserver

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yvanz/gin-tmpl/pkg/logger"
)

const (
	wrapUpTime = time.Second
	waitTime   = 5500 * time.Millisecond
)

var (
	wrapUpListeners          = new(listenerManager)
	shutdownListeners        = new(listenerManager)
	delayTimeBeforeForceQuit = waitTime
)

// AddShutdownListener adds fn as a shutdown listener.
// The returned func can be used to wait for fn getting called.
func AddShutdownListener(fn func()) (waitForCalled func()) {
	return shutdownListeners.addListener(fn)
}

// AddWrapUpListener adds fn as a wrap up listener.
// The returned func can be used to wait for fn getting called.
func AddWrapUpListener(fn func()) (waitForCalled func()) {
	return wrapUpListeners.addListener(fn)
}

func gracefulStop(signals chan os.Signal) {
	signal.Stop(signals)

	logger.Info("got signal SIGTERM, shutting down...")
	wrapUpListeners.notifyListeners()

	time.Sleep(wrapUpTime)
	shutdownListeners.notifyListeners()

	time.Sleep(delayTimeBeforeForceQuit - wrapUpTime)
	logger.Infof("still alive after %v, going to force kill the process...", delayTimeBeforeForceQuit)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}

type listenerManager struct {
	lock      sync.Mutex
	waitGroup sync.WaitGroup
	listeners []func()
}

func (lm *listenerManager) addListener(fn func()) (waitForCalled func()) {
	lm.waitGroup.Add(1)

	lm.lock.Lock()
	lm.listeners = append(lm.listeners, func() {
		defer lm.waitGroup.Done()
		fn()
	})
	lm.lock.Unlock()

	return func() {
		lm.waitGroup.Wait()
	}
}

func (lm *listenerManager) notifyListeners() {
	lm.lock.Lock()
	defer lm.lock.Unlock()

	for _, listener := range lm.listeners {
		listener()
	}
}
