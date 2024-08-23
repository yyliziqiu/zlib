package zhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/yyliziqiu/zlib/zutil"
)

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

func (cli *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, fkey string, fpath string, out interface{}) error {
	files := map[string]string{fkey: fpath}
	return cli.PostFormData(path, query, header, values, files, out)
}

func (cli *Client) PostStreamFormURL(path string, query url.Values, header http.Header, values map[string]string, filekey string, url string, out interface{}) error {
	data, _, err := cli.GetBinary(url, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostStream(path, query, header, values, filekey, filepath.Base(url), bytes.NewReader(data), out)
}

func (cli *Client) PostBinaryFormURL(path string, query url.Values, header http.Header, url string, out interface{}) error {
	data, typ, err := cli.GetBinary(url, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostBinary(path, query, header, typ, bytes.NewReader(data), out)
}
