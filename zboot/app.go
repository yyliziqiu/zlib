package zboot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyliziqiu/zlib/zconfig"
	"github.com/yyliziqiu/zlib/zlog"
)

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 配置文件路径
	ConfigPath string

	// 日志目录路径
	LogdirPath string

	// 全局配置
	Config any

	// 模块
	InitFuncs   InitFuncs
	BootFuncs   BootFuncs
	InitFuncsCb func() InitFuncs
	BootFuncsCb func() BootFuncs

	hasCallInitFuncs bool
}

// Init app
func (app *App) Init() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.CallInitFuncs()
}

func (app *App) InitConfig() (err error) {
	// 加载配置文件
	err = zconfig.Init(app.ConfigPath, app.Config)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	// 检查配置是否正确
	icheck, ok := app.Config.(ICheck)
	if ok {
		err = icheck.Check()
		if err != nil {
			return err
		}
	}

	// 为配置项设置默认值
	idefault, ok := app.Config.(IDefault)
	if ok {
		idefault.Default()
	}

	// 初始化日志
	logc := zlog.Config{}
	ilog, ok := app.Config.(IGetLog)
	if ok {
		logc = ilog.GetLog()
	}
	if app.LogdirPath != "" {
		logc.Path = app.LogdirPath
	}
	err = zlog.Init(logc)
	if err != nil {
		return fmt.Errorf("init log error [%v]", err)
	}

	return nil
}

func (app *App) CallInitFuncs() (err error) {
	if app.hasCallInitFuncs {
		return nil
	}
	app.hasCallInitFuncs = true

	initFuncs := app.InitFuncs
	if app.InitFuncsCb != nil {
		initFuncs = app.InitFuncsCb()
	}

	zlog.Info("Prepare init funcs.")
	err = initFuncs.Init()
	if err != nil {
		zlog.Errorf("Init funcs failed, error: %v", err)
		return err
	}
	zlog.Info("Init funcs succeed.")

	return nil
}

// Start app
func (app *App) Start() (err error, f context.CancelFunc) {
	err = app.InitConfig()
	if err != nil {
		return err, nil
	}
	return app.CallBootFuncs()
}

func (app *App) CallBootFuncs() (error, context.CancelFunc) {
	err := app.CallInitFuncs()
	if err != nil {
		return err, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	bootFuncs := app.BootFuncs
	if app.BootFuncsCb != nil {
		bootFuncs = app.BootFuncsCb()
	}

	zlog.Info("Prepare boot funcs.")
	err = bootFuncs.Boot(ctx)
	if err != nil {
		zlog.Errorf("Boot funcs failed, error: %v", err)
		cancel()
		return err, nil
	}
	zlog.Info("Boot funcs successfully.")

	return nil, cancel
}

// Run app
func (app *App) Run() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.CallBootFuncsBlocked()
}

func (app *App) CallBootFuncsBlocked() (err error) {
	err, cancel := app.CallBootFuncs()
	if err != nil {
		return err
	}

	zlog.Info("App run successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	zlog.Info("App prepare exit.")

	cancel()

	iwaittime, ok := app.Config.(IGetWaitTime)
	if ok {
		time.Sleep(iwaittime.GetWaitTime())
	}

	zlog.Info("App exit.")

	return nil
}
