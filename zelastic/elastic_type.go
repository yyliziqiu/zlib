package zelastic

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/zlib/zlog"
)

const DefaultId = "default"

type Config struct {
	Id           string   // optional
	Hosts        []string // must
	Username     string   // must
	Password     string   // must
	EnableLogger bool     // optional

	Logger *logrus.Logger // optional
	Client elastic.Doer   // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.EnableLogger && c.Logger == nil {
		c.Logger = zlog.NewWithNameMust("elastic-" + c.Id)
	}
	return c
}
