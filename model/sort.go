package model

import (
	"net/url"
	"fmt"
)

type Sort struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_sort",
		columns: []string{"uid", "name", "explain", "share", "first", "date", "editor", "status"},
		primary: "uid",
	}
	s := &Sort{m}
	models["Sort"] = s
}

func GetSort(v url.Values) *Sort {
	m := models["Sort"].(*Sort)
	m.values = v
	return m
}

// 路由
func (s *Sort) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Sort.Routing'key arg is nil", nil)
	}
	s.values = v
	m := map[string]func() []byte{
		"GetWithUid":    s.GetWithUid,
		"GetWithUids":   s.GetWithUids,
		"GetWithUser":   s.GetWithUser,
		"SetData":       s.SetDataWithPrimary,
		"GetData":       s.GetDataWithPrimary,
		"DelData":       s.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Sort is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条群组信息
func (s *Sort) GetWithUid() [] byte {
	uid := s.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Sort.GetWithSid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", s.getColumn(), s.table);
	return s.GetData(query, uid)
}

// 获取多条群组信息
func (s *Sort) GetWithUids() [] byte {
	uids := s.values.Get("uids")
	if len(uids) == 0 {
		return Result(0, "Sort.GetWithSids'uids is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` in (?)", s.getColumn(), s.table);
	return s.GetData(query, uids)
}

// 通过用户获取多条群组信息
func (s *Sort) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Sort.GetWithUser'user is nil", nil)
	}
	date := s.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE date > ? AND `uid` in (SELECT uid FROM jyh_join WHERE `user` = ?)", s.getColumn(), s.table)
	return s.GetData(query, date, user)
}

