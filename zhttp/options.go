package zhttp

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/sirupsen/logrus"
)

type Option func(cli *Client)

func WithClient(client *http.Client) Option {
	return func(cli *Client) {
		cli.client = client
	}
}

func WithCookie(o *cookiejar.Options) Option {
	return func(cli *Client) {
		jar, _ := cookiejar.New(o)
		cli.client.Jar = jar
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(cli *Client) {
		cli.client.Timeout = timeout
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(cli *Client) {
		cli.logger = logger
	}
}

func WithFormat(format string) Option {
	return func(cli *Client) {
		cli.format = format
	}
}

func WithError(error error) Option {
	return func(cli *Client) {
		cli.error = error
	}
}

func WithDumps(enabled bool) Option {
	return func(cli *Client) {
		cli.dumps = enabled
	}
}

func WithBaseURL(baseURL string) Option {
	return func(cli *Client) {
		cli.baseURL = baseURL
	}
}

func WithLogLength(n int) Option {
	return func(cli *Client) {
		cli.logLength = n
	}
}

func WithLogEscape(enabled bool) Option {
	return func(cli *Client) {
		cli.logEscape = enabled
	}
}

func WithRequestBefore(f func(r *http.Request)) Option {
	return func(cli *Client) {
		cli.requestBefore = f
	}
}

func WithBasicAuth(username string, password string) Option {
	return func(cli *Client) {
		cli.requestBefore = func(req *http.Request) {
			req.SetBasicAuth(username, password)
		}
	}
}

func WithBearerToken(token string) Option {
	return func(cli *Client) {
		cli.requestBefore = func(req *http.Request) {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func WithResponseAfter(f func(res *http.Response) error) Option {
	return func(cli *Client) {
		cli.responseAfter = f
	}
}

func Cookie(o *cookiejar.Options) Option {
	return WithCookie(o)
}

func Timeout(timeout time.Duration) Option {
	return WithTimeout(timeout)
}

func Logger(logger *logrus.Logger) Option {
	return WithLogger(logger)
}

func Format(format string) Option {
	return WithFormat(format)
}

func Error(error error) Option {
	return WithError(error)
}

func Dumps(enabled bool) Option {
	return WithDumps(enabled)
}

func BaseURL(baseURL string) Option {
	return WithBaseURL(baseURL)
}

func LogLength(n int) Option {
	return WithLogLength(n)
}

func LogEscape(enabled bool) Option {
	return WithLogEscape(enabled)
}

func RequestBefore(f func(r *http.Request)) Option {
	return WithRequestBefore(f)
}

func BasicAuth(username string, password string) Option {
	return WithBasicAuth(username, password)
}

func BearerToken(token string) Option {
	return WithBearerToken(token)
}

func ResponseAfter(f func(res *http.Response) error) Option {
	return WithResponseAfter(f)
}
