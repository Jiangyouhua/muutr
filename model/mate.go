package model

import (
	"net/url"
	"fmt"
)

type Mate struct {
	Model
}

func init() {
	initModels()
	ma := Model{
		table:   "jyh_mate",
		columns: []string{"uid", "sid", "user", "username", "role", "first", "date", "editor", "status"},
		primary: "uid",
	}
	m := &Mate{ma}
	models["Mate"] = m
}

func GetMate(v url.Values)*Mate{
	m := models["Mate"].(*Mate)
	m.values = v
	return m
}

// 路由
func (m *Mate) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Mate.Routing'key arg is nil", nil)
	}
	m.values = v
	ma := map[string]func() []byte{
		"GetWithUid":           m.GetWithUid,
		"GetWithSid":           m.GetWithSid,
		"GetWithUser":          m.GetWithUser,
		"SetUsernameWithUser": m.SetUsernameWithUser,
		"SetData":              m.SetDataWithPrimary,
		"GetData":              m.GetDataWithPrimary,
		"DelData":              m.DelDataWithPrimary,
	}
	f, ok := ma[k]
	if !ok {
		return Result(0, fmt.Sprintf("Mate is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条群组内某伙伴信息
func (m *Mate) GetWithUid() [] byte {
	uid := m.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Mate.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", m.getColumn(), m.table)
	return m.GetData(query, uid)
}

// 获取某群组所有伙伴信息
func (m *Mate) GetWithSid() [] byte {
	sid := m.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Mate.GetWithSid'sid is nil", nil)
	}
	date := m.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `sid` = ? AND date > ?", m.getColumn(), m.table)
	return m.GetData(query, sid, date)
}

// 通过用户获取多条群组信息
func (m *Mate) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Mate.GetWithUser'user is nil", nil)
	}
	date := m.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE date > ? AND `sid` in (SELECT sid FROM jyh_join WHERE `user` = ? )", m.getColumn(), m.table)
	return m.GetData(query, date, user)
}

// 修改用户在该群组的名称
func (m *Mate) SetUsernameWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Item.SetUsernameWithUser'user is nil", nil)
	}
	username := m.values.Get("username")
	if len(username) == 0 {
		return Result(0, "Item.SetUsernameWithUser'username is nil", nil)
	}
	query := fmt.Sprintf("UPDATE %s SET `username` = ? WHERE `user` = ?", m.table)
	return m.GetData(query, username, user)
}
