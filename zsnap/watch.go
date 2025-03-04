package zsnap

import (
	"context"
	"reflect"
	"time"

	"github.com/yyliziqiu/zlib/zlog"
	"github.com/yyliziqiu/zlib/zutil"
)

type Handler interface {
	Load() error
	Save() error
}

type HandlerName interface {
	Name() string
}

func handlerName(handler Handler) string {
	i, ok := handler.(HandlerName)
	if ok {
		return i.Name()
	}
	return reflect.TypeOf(handler).Name()
}

type HandlerInterval interface {
	Interval() time.Duration
}

func handlerInterval(handler Handler) time.Duration {
	i, ok := handler.(HandlerInterval)
	if ok {
		return i.Interval()
	}
	return 0
}

func Watch(ctx context.Context, handlersFunc func() []Handler) error {
	handlers := handlersFunc()

	err := watchLoad(handlers)
	if err != nil {
		return err
	}

	for _, handler := range handlers {
		go runWatchSave(ctx, handler)
	}

	return nil
}

func watchLoad(handlers []Handler) error {
	timer := zutil.NewTimer()
	for _, handler := range handlers {
		err := handler.Load()
		if err != nil {
			zlog.Errorf("Load snap failed, name: %s, error: %v.", handlerName(handler), err)
			return err
		}
		zlog.Infof("Load snap succeed, name: %s, cost: %s.", handlerName(handler), timer.Pauses())
	}
	zlog.Infof("Load snaps compeleted, cost: %s.", timer.Stops())
	return nil
}

func runWatchSave(ctx context.Context, handler Handler) {
	interval := handlerInterval(handler)

	if interval <= 0 {
		<-ctx.Done()
		_ = watchSave(handler)
		return
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			_ = watchSave(handler)
		case <-ctx.Done():
			_ = watchSave(handler)
			return
		}
	}
}

func watchSave(handler Handler) error {
	timer := zutil.NewTimer()
	err := handler.Save()
	if err != nil {
		zlog.Errorf("Save snap failed, name: %s, error: %v.", handlerName(handler), err)
	} else {
		zlog.Infof("Save snap succeed, name: %s, cost: %s.", handlerName(handler), timer.Stops())
	}
	return err
}
