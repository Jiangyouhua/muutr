package model

import (
	"fmt"
	"net/url"
)

// 聊天表
type Config struct {
	Model
}

// 初始表信息
func init() {
	initModels()
	m := &Model{
		table:   "jyh_config",
		columns: []string{"uid", "server", "back", "zip", "file", "upload", "socket", "html", "date"},
		primary: "uid",
	}
	c := &Config{*m}
	models["Config"] = c
}

func GetConfig(v url.Values)*Config{
	m := models["Config"].(*Config)
	m.values = v
	return m
}

// 路由
func (c *Config) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Config.Routing'key arg is nil", nil)
	}
	c.values = v
	m := map[string]func() []byte{
		"SetData":          c.SetConfig,
		"GetData":          c.GetConfig,
		"DelData":          c.DelConfig,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Config is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条聊天数据
func (c *Config) SetConfig() []byte {
	c.values.Set("uid", "1");
	return c.SetDataWithPrimary();
}

// 获取单条聊天数据
func (c *Config) GetConfig() []byte {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = 1", c.getColumn(), c.table)
	return c.Query(query);
}

func (c *Config)DelConfig()[]byte {
	return Result(0, "Config can't Del", nil)
}
