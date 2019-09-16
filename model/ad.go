/**
 * 需要有效性，
 * 1. APP在请求之前，保证本地所有的无效
 * 2. 请求最新一个。如果有，则同步到本地，供APP显示；如果没有，则APP显示默认页面
 */
package model

import (
	"fmt"
	"net/url"
)

type Ad struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_ad",
		columns: []string{"uid", "image", "link", "start", "end", "date", "admin", "status", "key"},
		primary: "uid",
	}
	a := &Ad{m}
	models["Ad"] = a
}

func GetAd(v url.Values)*Ad{
	m := models["Ad"].(*Ad)
	m.values = v
	return m
}

// 路由
func (a *Ad) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "User.Routing'key arg is nil", nil)
	}
	a.values = v
	m := map[string]func() []byte{
		"GetLatest": a.GetLatest,
		"GetAll":    a.GetAll,
		"SetData":   a.SetDataByAdmin,
		"GetData":   a.GetDataWithPrimary,
		"DelData":   a.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("ad is not this func (%s)", k), nil)
	}
	return f()
}

// 获取最新一个广告信息，APP上会排他处理，需要有效性
func (a *Ad) GetLatest() [] byte {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `status` = 1 AND `start` < NOW() ORDER BY `start` DESC LIMIT 1", a.getColumn(), a.table)
	return a.GetData(query)
}

// 获取全部
func (a *Ad) GetAll() []byte {
	return a.Model.GetDataWithPage(a.values.Get("page"), "", "`start` DESC")
}

func(a *Ad) SetDataByAdmin()[]byte{
	a.values.Set("admin", Session.Get("user").(string))
	return a.SetDataWithPrimary()
}
