package model

import (
	"net/url"
	"fmt"
)

type Item struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_item",
		columns: []string{"uid", "sid", "explain", "cycle", "option", "start", "end",  "complete", "timeline", "first", "date", "rank", "mine", "editor", "status"},
		primary: "uid",
	}
	i := &Item{m}
	models["Item"] = i
}

func GetItem(v url.Values)*Item{
	m := models["Item"].(*Item)
	m.values = v
	return m
}

// 路由
func (i *Item) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Item.Routing'key arg is nil", nil)
	}
	i.values = v
	m := map[string]func() []byte{
		"GetWithUid":          i.GetWithUid,
		"GetWithSid":          i.GetWithSid,
		"GetWithUser":         i.GetWithUser,
		"SetData":             i.SetDataWithPrimary,
		"GetData":             i.GetDataWithPrimary,
		"DelData":             i.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Item is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条目标信息
func (i *Item) GetWithUid() [] byte {
	uid := i.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Item.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", i.getColumn(), i.table)
	return i.GetData(query, uid)
}

// 获取群组内，最新目标信息
func (i *Item) GetWithSid() [] byte {
	sid := i.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Item.GetWithSid'sid is nil", nil)
	}
	date := i.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `sid` = ? AND `date` > ? ORDER BY uid DESC", i.getColumn(), i.table)
	return i.GetData(query, sid, date)
}

func (i *Item) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Item.GetWithUser'user is nil", nil)
	}
	date := i.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE date > ? AND `sid` in (SELECT sid FROM jyh_join WHERE user = ? AND status = 1)  ORDER BY uid DESC", i.getColumn(), i.table)
	return i.GetData(query, date, user)
}
