package zweb

import (
	"github.com/sirupsen/logrus"
)

var (
	_accessLogger *logrus.Logger
	_errorLogger  *logrus.Logger
)

type Config struct {
	Addr             string
	DisableAccessLog bool
	ErrorLogName     string
	AccessLogName    string
}

func (c Config) Default() Config {
	if c.Addr == "" {
		c.Addr = ":80"
	}
	if c.ErrorLogName == "" {
		c.ErrorLogName = "web-error"
	}
	if c.AccessLogName == "" {
		c.AccessLogName = "web-access"
	}
	return c
}
