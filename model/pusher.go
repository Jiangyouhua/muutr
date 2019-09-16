/**
 * 需要有效性，
 * 1. 推送信息一个接一个的请求，推送后置为无效
 * 2. 请求最新一个。如果有，则同步到服务器内存，按指定的时间推送。
 */
package model

import (
	"fmt"
	"net/url"
)

type Pusher struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_pusher",
		columns: []string{"uid", "title", "body", "image", "category", "push", "admin", "status"},
		primary: "uid",
	}
	p := &Pusher{m}
	models["Pusher"] = p
}

func GetPusher(v url.Values)*Pusher{
	m := models["Pusher"].(*Pusher)
	m.values = v
	return m
}

// 路由
func (p *Pusher) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "User.Routing'key arg is nil", nil)
	}
	p.values = v
	m := map[string]func() []byte{
		"GetAll":				p.GetAll,
		"SetData":              p.SetDataByAdmin,
		"GetData":              p.GetDataWithPrimary,
		"DelData":              p.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("pusher is not this func (%s)", k), nil)
	}
	return f()
}

// 获取全部
func (p *Pusher)GetAll()[]byte{
	return p.Model.GetDataWithPage(p.values.Get("page"), "", "`push` DESC")
}

// 插入的同时更新最近一条需要发布的信息
func (p *Pusher)SetDataByAdmin()[]byte{
	// 更新
	p.values.Set("admin", Session.Get("user").(string))
	b := p.SetDataWithPrimary()
	go pushing()
	return b
}


// 获取最新一条推送
func (p *Pusher)GetLatestOne()(data map[string]string , err error){
	// 没有更新就不要读数据库
	query := fmt.Sprintf("SELECT `title`, `body`, `image`, `category`, `push` FROM %s WHERE `status` = 1 AND `push` >= NOW()  ORDER BY `push` LIMIT 1",  p.table)
	re, err := p.Fetch(query)
	if err != nil {
		return
	}
	data = re[0]
	return
}

