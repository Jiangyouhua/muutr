package model

import (
	"fmt"
	"net/url"
)

type Record struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_record",
		columns: []string{"uid", "cid", "user", "date", "editor", "status"},
		primary: "uid",
	}
	u := &Record{m}
	models["Record"] = u
}

func GetRecord(v url.Values) *Record {
	m := models["Record"].(*Record)
	m.values = v
	return m
}

// 路由
func (r *Record) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Record.Routing'key arg is nil", nil)
	}
	r.values = v
	m := map[string]func() []byte{
		"StatusWithCidWithUser": r.StatusWithCidWithUser,
		"SetData":               r.SetDataWithPrimary,
		"GetData":               r.GetDataWithPrimary,
		"DelData":               r.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Router is not this func (%s)", k), nil)
	}
	return f()
}

// 保存修改
func (r *Record) StatusWithCidWithUser() []byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Record.StatusWithCidWithUser'cid is nil", nil)
	}
	cid := r.values.Get("cid")
	if len(user) == 0 {
		return Result(0, "Record.StatusWithCidWithUser'user is nil", nil)
	}
	return r.Exce("UPDATE jyh_record SET `status` = 0 WHERE `cid` = ? AND `user` = ?", cid, user)
}
