package zhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (cli *Client) logHTTP(log HTTPLog) {
	if cli.logger == nil {
		return
	}

	headers, reqBody, resBody := SerializeHeader(log.Request.Header), "", ""
	if len(log.RequestBody) > 0 {
		reqBody = string(log.RequestBody)
	}
	if len(log.ResponseBody) > 0 {
		resBody = string(log.ResponseBody)
	}

	if log.Error == nil {
		cli.logInfo("Response(%d Succeed), URL: %s, headers: %s, request: %s, response: %s, cost: %s.",
			log.Response.StatusCode, log.Request.URL, headers, reqBody, resBody, log.Cost)
	} else {
		cli.logWarn("Response(%d Failed), URL: %s, headers: %s, request: %s, response: %s, error: %v, cost: %s.",
			log.Response.StatusCode, log.Request.URL, headers, reqBody, resBody, log.Error, log.Cost)
	}
}

func (cli *Client) logInfo(format string, args ...interface{}) {
	if cli.logger == nil {
		return
	}
	message := cli.logCheck(fmt.Sprintf(format, args...))
	cli.logger.Info(message)
}

func (cli *Client) logWarn(format string, args ...interface{}) {
	if cli.logger == nil {
		return
	}
	message := cli.logCheck(fmt.Sprintf(format, args...))
	cli.logger.Warn(message)
}

func (cli *Client) logCheck(log string) string {
	if cli.logLength <= 0 {
		return ""
	}

	if len(log) > cli.logLength {
		log = log[:cli.logLength]
	}

	if cli.logEscape {
		log = strings.ReplaceAll(log, "\t", "\\t")
		log = strings.ReplaceAll(log, "\r", "\\r")
		log = strings.ReplaceAll(log, "\n", "\\n")
	}

	return log
}

func (cli *Client) dumpRequest(req *http.Request) {
	if !cli.dumps {
		return
	}

	bs, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Dump request failed, error: %v\n", err)
		return
	}

	fmt.Println("\n---------- Request ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Request End----------")
}

func (cli *Client) dumpResponse(res *http.Response) {
	if !cli.dumps {
		return
	}

	bs, err := httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Printf("Dump response failed, error: %v", err)
		return
	}

	fmt.Println("\n---------- Response ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Response End----------")
}
