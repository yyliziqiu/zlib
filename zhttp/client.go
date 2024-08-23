package zhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/zlib/zutil"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

type Client struct {
	client        *http.Client                   //
	logger        *logrus.Logger                 // 如果为 nil，则不记录日志
	format        string                         //
	error         error                          // 不能是指针
	dumps         bool                           // 将 HTTP 报文打印到控制台
	baseURL       string                         //
	logLength     int                            // 日志最大长度
	logEscape     bool                           // 替换日志中的特殊字符
	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *Client {
	client := &Client{
		client:        &http.Client{Timeout: 5 * time.Second},
		logger:        nil,
		format:        FormatJSON,
		error:         nil,
		dumps:         false,
		baseURL:       "",
		logLength:     1024,
		logEscape:     false,
		requestBefore: nil,
		responseAfter: nil,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (cli *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinURL(cli.baseURL, path)
	}

	url2, err := AppendQuery(path, query)
	if err != nil {
		cli.logWarn("Append query failed, URL: %s, query: %s, error: %v.", url2, query.Encode(), err)
		return nil, fmt.Errorf("append query error [%v]", err)
	}

	req, err := http.NewRequest(method, url2, body)
	if err != nil {
		cli.logWarn("New request failed, URL: %s, error: %v.", url2, err)
		return nil, fmt.Errorf("new request error [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if cli.requestBefore != nil {
		cli.requestBefore(req)
	}

	return req, nil
}

func (cli *Client) doRequest(req *http.Request) (*http.Response, error) {
	cli.dumpRequest(req)

	res, err := cli.client.Do(req)
	if err != nil {
		cli.logWarn("Do request failed, URL: %s, error: %v.", req.URL, err)
		return nil, err
	}

	return res, nil
}

func (cli *Client) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	cli.dumpResponse(res)

	if cli.responseAfter != nil {
		err := cli.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error [%v]", err)
	}

	switch cli.format {
	case FormatText:
		return body, cli.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, cli.handleJSONResponse(res.StatusCode, body, out)
	}
}

func (cli *Client) handleJSONResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 == 2 {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					err2, ok2 := out.(error)
					if ok2 {
						return err2
					}
					return newResponseError(statusCode, string(body))
				}
			}
		}
	} else {
		if cli.error != nil {
			ret := reflect.New(reflect.TypeOf(cli.error)).Interface()
			err := json.Unmarshal(body, ret)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			return ret.(error)
		} else if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			err2, ok2 := out.(error)
			if ok2 {
				return err2
			}
			return newResponseError(statusCode, string(body))
		}
	}
	return nil
}

func (cli *Client) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 != 2 {
		return newResponseError(statusCode, string(body))
	}

	if out == nil {
		return nil
	}

	bs, ok := out.(*[]byte)
	if !ok {
		return fmt.Errorf("response receiver must *[]byte type")
	}
	*bs = body

	return nil
}

func (cli *Client) get(method string, path string, query url.Values, header http.Header, out interface{}) error {
	req, err := cli.newRequest(method, path, query, header, nil)
	if err != nil {
		return err
	}

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := cli.handleResponse(res, out)

	cli.logHTTP(HTTPLog{
		Method:       method,
		Request:      req,
		RequestBody:  nil,
		Response:     res,
		ResponseBody: body,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) post(method string, path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := cli.newRequest(method, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logHTTP(HTTPLog{
		Method:       method,
		Request:      req,
		RequestBody:  reqBody,
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) Get(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodGet, path, query, header, out)
}

func (cli *Client) Post(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPost, path, query, header, in, out)
}

func (cli *Client) Put(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPut, path, query, header, in, out)
}

func (cli *Client) Delete(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodDelete, path, query, header, out)
}

func (cli *Client) PostJSON(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.Post(path, query, header, in, out)
}

func (cli *Client) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()

	req, err := cli.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	reqBody, _ = url.QueryUnescape(reqBody)
	cli.logHTTP(HTTPLog{
		Method:       http.MethodPost,
		Request:      req,
		RequestBody:  []byte(reqBody),
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) PostFormData(path string, query url.Values, header http.Header, values map[string]string, files map[string]string, out interface{}) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
	}

	if len(files) > 0 {
		for key, file := range files {
			err := cli.writeFormFile(writer, key, file)
			if err != nil {
				return err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	var cpy map[string]string
	var bod []byte
	if len(values) > 0 {
		cpy = values
		for key, file := range files {
			cpy[key] = file
		}
	} else {
		cpy = files
	}
	if len(cpy) == 0 {
		bod = []byte("{}")
	} else {
		bod, _ = json.Marshal(cpy)
	}
	cli.logHTTP(HTTPLog{
		Method:       http.MethodPost,
		Request:      req,
		RequestBody:  bod,
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) writeFormFile(writer *multipart.Writer, key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(key, file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) GetBinary(path string, query url.Values, header http.Header) ([]byte, string, error) {
	req, err := cli.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return nil, "", err
	}

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)

	cli.logHTTP(HTTPLog{
		Method:       http.MethodGet,
		Request:      req,
		RequestBody:  nil,
		Response:     res,
		ResponseBody: nil,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return bs, res.Header.Get("Content-Type"), err
}

func (cli *Client) PostBinary(path string, query url.Values, header http.Header, mimeType string, in io.Reader, out interface{}) error {
	req, err := cli.newRequest(http.MethodPost, path, query, header, in)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mimeType)

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logHTTP(HTTPLog{
		Method:       http.MethodPost,
		Request:      req,
		RequestBody:  nil,
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, filekey string, filepath string, out interface{}) error {
	files := map[string]string{filekey: filepath}
	return cli.PostFormData(path, query, header, values, files, out)
}

func (cli *Client) PostStream(path string, query url.Values, header http.Header, values map[string]string, filekey string, filename string, stream io.Reader, out interface{}) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
	}

	if stream != nil {
		part, err := writer.CreateFormFile(filekey, filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, stream)
		if err != nil {
			return err
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	timer := zutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	var bod []byte
	if len(values) == 0 {
		bod = []byte(fmt.Sprintf(`{"%s":"%s"}`, filekey, filename))
	} else {
		values[filekey] = filename
		bod, _ = json.Marshal(values)
	}
	cli.logHTTP(HTTPLog{
		Method:       http.MethodPost,
		Request:      req,
		RequestBody:  bod,
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}
