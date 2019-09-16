package model

import (
	"net/url"
	"fmt"
)

type Token struct {
	Model
}

// 初始表信息
func init() {
	initModels()
	m := &Model{
		table:   "jyh_token",
		columns: []string{"token", "ip", "user", "status"},
		primary: "token",
	}
	c := &Token{*m}
	models["Token"] = c
}

func GetToken(v url.Values)*Token{
	m := models["Token"].(*Token)
	m.values = v
	return m
}

// 路由
func (t *Token) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Token.Routing'key arg is nil", nil)
	}
	t.values = v
	m := map[string]func() []byte{
		"SetData": t.SetDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Token is not this func (%s)", k), nil)
	}
	return f()
}

func (t *Token)GetAllToken()([]map[string]string, error){
	t.values.Set("ip", ToString(Session.Get("ip")))
	return t.Fetch("SELECT `token` FROM jyh_token")
}

func(t *Token)SetData()[]byte{
	return  t.SetDataWithPrimary()
}

