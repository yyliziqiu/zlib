package zboot

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/zlib/zdb"
	"github.com/yyliziqiu/zlib/zelastic"
	"github.com/yyliziqiu/zlib/zenv"
	"github.com/yyliziqiu/zlib/zkafka"
	"github.com/yyliziqiu/zlib/zlog"
	"github.com/yyliziqiu/zlib/zredis"
	"github.com/yyliziqiu/zlib/ztask"
	"github.com/yyliziqiu/zlib/zweb"
)

// ICheck 检查配置是否正确
type ICheck interface {
	Check() error
}

// IDefault 为配置项设置默认值
type IDefault interface {
	Default()
}

// IGetLog 获取日志配置
type IGetLog interface {
	GetLog() zlog.Config
}

// IGetWaitTime 获取应用退出时等待时长配置
type IGetWaitTime interface {
	GetWaitTime() time.Duration
}

type Config struct {
	Env      string
	AppId    string
	InsId    string
	BasePath string
	DataPath string
	WaitTime time.Duration

	Log zlog.Config
	Web zweb.Config

	DB      []zdb.Config
	Redis   []zredis.Config
	Kafka   []zkafka.Config
	Elastic []zelastic.Config

	Migration struct {
		EnableTables  bool
		EnableRecords bool
	}

	CronTask []ztask.CronTask
	OnceTask []ztask.OnceTask

	Values map[string]string
}

func (c *Config) Default() {
	if c.Env == "" {
		c.Env = zenv.Prod
	}
	if c.AppId == "" {
		c.AppId = "app"
	}
	if c.InsId == "" {
		c.InsId = "1"
	}
	if c.BasePath == "" {
		c.BasePath = "."
	}
	if c.DataPath == "" {
		c.DataPath = filepath.Join(c.BasePath, "data")
	}
}

func (c *Config) GetLog() zlog.Config {
	return c.Log
}

func (c *Config) GetWaitTime() time.Duration {
	return c.WaitTime
}
