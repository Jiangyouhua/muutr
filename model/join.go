package model

import (
	"net/url"
	"fmt"
	"strconv"
	"encoding/json"
)

type Join struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_join",
		columns: []string{"uid", "sid", "name", "user", "sequence", "first", "date", "editor", "status"},
		primary: "uid",
	}
	j := &Join{m}
	models["Join"] = j
}

func GetJoin(v url.Values) *Join {
	m := models["Join"].(*Join)
	m.values = v
	return m
}

// 路由
func (j *Join) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Join.Routing'key arg is nil", nil)
	}
	j.values = v
	m := map[string]func() []byte{
		"GetWithUid":  j.GetWithUid,
		"GetWithUser": j.GetWithUser,
		"SetRole":     j.SetRole,
		"SetSequence": j.SetSequence,
		"SetData":     j.SetDataWithPrimary,
		"GetData":     j.GetDataWithPrimary,
		"DelData":     j.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Join is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条用户加入群组信息
func (j *Join) GetWithUid() [] byte {
	uid := j.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Join.GetWithUid.uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", j.getColumn(), j.table)
	return j.GetData(query, uid)
}

// 获取用户加入的群组信息
func (j *Join) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Join.GetWithUser'user is nil", nil)
	}
	date := j.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE AND `user` = ? AND date > ?", j.getColumn(), j.table)
	return j.GetData(query, user, date)
}

func (j *Join) SetRole() []byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Join.SetRole'user is nil", nil)
	}
	role := j.values.Get("role")
	i, err := strconv.Atoi(role)
	if err != nil {
		return Result(0, err.Error(), nil)
	}
	adminRole := Session.Get("role").(int)
	if adminRole <= i {
		return Result(0, "Beyond the Limits of your authority", nil)
	}
	query := fmt.Sprintf("INSERT INTO %s (`user`, `role`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `role` = VALUES(`role`)", j.table)
	return j.SetData(query, user, role)
}

func (j *Join) SetSequence() []byte {
	data := j.values.Get("data")
	if len(data) == 0 {
		return Result(0, "Join.SetSequence'data is nil", nil)
	}
	var i interface{}
	err := json.Unmarshal([]byte(data), &i)
	if err != nil {
		return Result(0, "Join.SetSequence'data json.Unmarshal is err", nil)
	}
	query := j.InsertUpdateStr(i)
	return j.Exce(query)
}

