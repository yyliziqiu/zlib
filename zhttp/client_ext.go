package zhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"

	"github.com/yyliziqiu/zlib/zmime"
	"github.com/yyliziqiu/zlib/zutil"
)

func (cli *Client) postStream(
	path string,
	query url.Values,
	header http.Header,
	values map[string]string,
	field string,
	filename string,
	mimeType string,
	stream io.Reader,
	out interface{},
) error {
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
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, EscapeQuotes(field), EscapeQuotes(filename)))
		h.Set("Content-Type", mimeType)
		part, err := writer.CreatePart(h)
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
		bod = []byte(fmt.Sprintf(`{"%s":"%s"}`, field, filename))
	} else {
		values[field] = filename
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

func (cli *Client) PostStream(path string, query url.Values, header http.Header, values map[string]string, field string, filename string, stream io.Reader, out interface{}) error {
	mimeType := zmime.Get(filename)
	return cli.postStream(path, query, header, values, field, filename, mimeType, stream, out)
}

func (cli *Client) PostStream2(path string, query url.Values, header http.Header, values map[string]string, field string, filename string, stream io.Reader, out interface{}) error {
	mimeType := "application/octet-stream"
	return cli.postStream(path, query, header, values, field, filename, mimeType, stream, out)
}

func (cli *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, field string, filepath string, out interface{}) error {
	files := map[string]string{field: filepath}
	return cli.PostFormData(path, query, header, values, files, out)
}

func (cli *Client) PostBinaryFormURL(path string, query url.Values, header http.Header, url string, out interface{}) error {
	data, typ, err := cli.GetBinary(url, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostBinary(path, query, header, typ, bytes.NewReader(data), out)
}

func (cli *Client) PostStreamFormURL(path string, query url.Values, header http.Header, values map[string]string, field string, url string, out interface{}) error {
	data, _, err := cli.GetBinary(url, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostStream(path, query, header, values, field, filepath.Base(url), bytes.NewReader(data), out)
}

func (cli *Client) PostStreamFormURL2(path string, query url.Values, header http.Header, values map[string]string, field string, url string, out interface{}) error {
	data, _, err := cli.GetBinary(url, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostStream2(path, query, header, values, field, filepath.Base(url), bytes.NewReader(data), out)
}
