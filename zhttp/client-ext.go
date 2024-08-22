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

	"github.com/yyliziqiu/zlib/zutil"
)

func (cli *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, files map[string]string, out interface{}) error {
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

func (cli *Client) Post(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.PostJSON(path, query, header, in, out)
}

func (cli *Client) Put(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := cli.newRequest(http.MethodPut, path, query, header, bytes.NewReader(reqBody))
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
		Request:      req,
		RequestBody:  reqBody,
		Response:     res,
		ResponseBody: resBody,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}

func (cli *Client) Delete(path string, query url.Values, header http.Header, out interface{}) error {
	req, err := cli.newRequest(http.MethodDelete, path, query, header, nil)
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
		Request:      req,
		RequestBody:  nil,
		Response:     res,
		ResponseBody: body,
		Error:        err,
		Cost:         timer.Stops(),
	})

	return err
}
