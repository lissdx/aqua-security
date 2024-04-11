package fswatchdog

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lissdx/aqua-security/internal/pkg/consumer"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/spf13/viper"
	"sync"
	"sync/atomic"
)

var _ consumer.Consumer = (*FsWatchdog)(nil)

const watchingDirConfigParamName = "WATCHING_DIR_PATH"

type FsWatchdog struct {
	fsWatcher    *fsnotify.Watcher
	activeStatus atomic.Value
	logger       asLogger.TOLogger
	doneChannel  chan interface{}
	allProcessWg *sync.WaitGroup
	outMsgStream <-chan interface{}
	mu           sync.Mutex
}

func (fsw *FsWatchdog) Run() {
	msgStream := make(chan interface{})
	//fsw.doneChannel = make(chan interface{})

	fsw.outMsgStream = msgStream
	fsw.activeStatus.Store(uint8(1))
	wg := sync.WaitGroup{} // wait for all subprocesses ended and then close delete stream
	fsw.allProcessWg = &sync.WaitGroup{}

	//fsw.fsWatcher, err = fsnotify.NewWatcher()

	fsw.allProcessWg.Add(1)
	wg.Add(1)

	// main event send process
	go func(streamName string, wg *sync.WaitGroup) {
		defer wg.Done()

		fsw.logger.Info(fmt.Sprintf("FS Watchdog: process %s started", streamName))
		isWatching := true
		for isWatching {
			select {
			// Read from Errors.
			case err, ok := <-fsw.fsWatcher.Errors:
				if !ok { // Channel was closed (i.e. Watcher.Close() was called).
					isWatching = false
				} else {
					fsw.logger.Error("FsWatcher ERROR: %s", err.Error())
				}
				//printTime("ERROR: %s", err)
			// Read from Events.
			case e, ok := <-fsw.fsWatcher.Events:
				if !ok { // Channel was closed (i.e. Watcher.Close() was called).
					isWatching = false
				} else {

					// We just want to watch for file creation, so ignore everything
					// outside of Create and Write.
					if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
						continue
					}

					msgStream <- e
				}
				// Just print the event nicely aligned, and keep track how many
				// events we've seen.
				//i++
				//printTime("%3d %s", i, e)
			}
		}

		fsw.logger.Info(fmt.Sprintf("Watchdog: watching process %s end", streamName))
	}(fmt.Sprintf("watchdog_proccess_%d", 0), &wg)

	// watching process
	go func(group *sync.WaitGroup, totalWait *sync.WaitGroup) {
		defer totalWait.Done()
		defer close(msgStream)
		group.Wait()
		fsw.logger.Info("Watchdog consumer: all processes stopped.")
	}(&wg, fsw.allProcessWg)

}

func (fsw *FsWatchdog) Stop() {
	defer fsw.mu.Unlock()
	fsw.mu.Lock()

	if fsw.isActive() {
		fsw.setActiveOff()
		_ = fsw.fsWatcher.Close()
		fsw.allProcessWg.Wait()
		fsw.logger.Info("Watchdog consumer service STOP")
	} else {
		fsw.logger.Info("Watchdog consumer service is not active")
	}
}

func (fsw *FsWatchdog) ConsumerStream() consumer.InStream {
	//TODO implement me
	return fsw.outMsgStream
}

func (fsw *FsWatchdog) isActive() bool {
	return fsw.activeStatus.Load() != nil && fsw.activeStatus.Load().(uint8) > 0
}

func (fsw *FsWatchdog) setActiveOff() {
	fsw.activeStatus.Store(uint8(0))
}

func NewFsWatchdog(vConf *viper.Viper, logger asLogger.TOLogger) consumer.Consumer {

	wDir := strUtil.NormalizeString(vConf.GetString(watchingDirConfigParamName))
	if strUtil.IsEmptyString(wDir) {
		panic(fmt.Sprintf("error on NewFsWatchdog: %s param is mandatory", watchingDirConfigParamName))
	}

	fsWatchdog := &FsWatchdog{
		logger: logger,
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		if w != nil {
			err = w.Close()
			if err != nil {
				logger.Error("error on NewFsWatchdog Close: %s", err.Error())
			}
		}
		panic(fmt.Sprintf("error on fsnotify.NewWatcher(): %s", err.Error()))
	}

	err = w.Add(wDir)
	if err != nil {
		_ = w.Close()
		panic(fmt.Sprintf("error on fsnotify.NewWatcher(): %s\n", err.Error()))
	}

	fsWatchdog.fsWatcher = w

	logger.Info(fmt.Sprintf("Watchdog is starting to listen dir: %s", wDir))
	return fsWatchdog
}
