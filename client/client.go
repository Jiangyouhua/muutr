package client

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
	"ido/logger"
	"fmt"
)

const (
	MethodGet = iota
	MethodPost
	MethodForm
)

var (
	Log *logger.Logger
)

// Client 客户端HTTP请求
type Client struct {
	URL  string  // 地址
	Method int  // 模式
	Data interface{}  // 数据
}

func (c *Client)Run()[]byte{
	var r *http.Response
	// 按不同的请求模式分别请求
	switch c.Method {
	case MethodGet:
		r = c.Get()
	case MethodPost:
		r = c.Post()
	case MethodForm:
		r = c.Form()
	}
	return c.Result(r)
}

// 获取二进制返回结果
func (c *Client) Result(r *http.Response) []byte {
	if r == nil {
		return nil
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	switch {
	case r.StatusCode >= 500:
		Log.Out(logger.ERR_LOG, fmt.Sprintf("HttpServer is err, %s", r.Status))
		return nil
	case r.StatusCode >= 400:
		Log.Out(logger.ERR_LOG, fmt.Sprintf("HttpRequest is bad, %s", r.Status))
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	return body
}

// Get 客户端Get请求
func (c *Client) Get() *http.Response {
	// fmt.Println(url)
	r, err := http.Get(c.URL)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	return r
}

// Post 客户端Post请求
func (c *Client) Post() *http.Response {
	j, err := json.Marshal(c.Data)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	s := string(j)
	r, err := http.Post(c.URL, "application/x-www-form-urlencoded", strings.NewReader(s))
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	return r
}

// Form 客户端Form请求
func (c *Client) Form() *http.Response {
	r, err := http.PostForm(c.URL, c.Data.(url.Values))
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	return r
}
