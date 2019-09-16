package model

import (
	"fmt"
	"net/url"
)

type Claim struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_claim",
		columns: []string{"uid", "sid", "user", "iid", "alert", "date", "editor", "status"},
		primary: "uid",
	}
	t := &Claim{m}
	models["Claim"] = t
}

func GetClaim(v url.Values)*Claim {
	m := models["Claim"].(*Claim)
	m.values = v
	return m
}

// 路由
func (c *Claim) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "User.Routing'key arg is nil", nil)
	}
	c.values = v
	m := map[string]func() []byte{
		"GetWithUid":  c.GetWithUid,
		"GetWithUser": c.GetWithUser,
		"SetData":     c.SetDataWithPrimary,
		"GetData":     c.GetDataWithPrimary,
		"DelData":     c.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("task is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单人用户信息
func (c *Claim) GetWithUid() [] byte {
	uid := c.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Claim.GetByUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", c.getColumn(), c.table)
	return c.GetData(query, uid)
}

func (c *Claim) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Claim.GetWithUser'user is nil", nil)
	}
	date := c.values.Get("date")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE date > ? AND `uid` in (SELECT uid FROM jyh_join WHERE `user` = ? AND `status` = 1)  ORDER BY iid DESC", c.getColumn(), c.table)
	return c.GetData(query, date, user)
}
