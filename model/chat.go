/**
 * 需要有效性，
 * 1. APP在请求之前，保证本地只保留有效的一页数（有效性供敏感词过滤用）
 * 2. 请求得到均为有效。
 */
package model

import (
	"fmt"
	"net/url"
	"strconv"
)

// 聊天表
type Chat struct {
	Model
}

// 初始表信息
func init() {
	initModels()
	m := &Model{
		table:   "jyh_chat",
		columns: []string{"uid", "sid", "user", "date", "content", "username", "at", "analysis", "status"},
		primary: "uid",
	}
	c := &Chat{*m}
	models["Chat"] = c
}

func GetChat(v url.Values) *Chat {
	m := models["Chat"].(*Chat)
	m.values = v
	return m
}

// 路由
func (c *Chat) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Chat.Routing'key arg is nil", nil)
	}
	c.values = v
	m := map[string]func() []byte{
		"GetWithUid":          c.GetWithUid,
		"GetWithSidStartDate": c.GetWithSidStartDate,
		"GetWithSidEndDate":   c.GetWithSidEndDate,
		"SetData":             c.SetData,
		"GetData":             c.GetDataWithPrimary,
		"DelData":             c.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Chat is not this func (%s)", k), nil)
	}
	return f()
}

func (c *Chat) whereWithAt() string {
	user, ok := Session.Get("user").(string)
	if !ok || len(user) == 0 {
		return ""
	}
	return fmt.Sprintf("AND (`user` = %s OR `at` = '' OR `at` IS NULL OR FIND_IN_SET(%s,`at`))", user, user)
}

// 获取单条聊天数据
func (c *Chat) GetWithUid() []byte {
	uid := c.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Chat.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", c.getColumn(), c.table);
	return c.GetData(query, uid)
}

// 获取群组内最新N条聊天数据
func (c *Chat) GetWithSidStartDate() []byte {
	sid := c.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Chat.GetWithSid'sid is nil", nil)
	}

	num := c.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "50"
	}

	date := c.values.Get("date")
	if len(date) == 0 {
		date = datetime + " 000"
	} else {
		num = "1000"
	}

	// 实现at
	where := c.whereWithAt()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `sid` = ? AND date > ? %s ORDER BY `date` DESC LIMIT ?", c.getColumn(), c.table, where)
	return c.GetData(query, sid, date, num)
}

// 分页，以CID为起点，获取群组内最新N条聊天数据
func (c *Chat) GetWithSidEndDate() []byte {
	sid := c.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Chat.GetWithSidEndDate'sid is nil", nil)
	}
	date := c.values.Get("date")
	if len(date) == 0 {
		return Result(0, "Chat.GetWithSidEndDate'date is nil", nil)
	}
	num := c.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "50"
	}

	// 实现at
	where := c.whereWithAt()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `sid` = ? AND date < ? %s ORDER BY `date` DESC LIMIT ?", c.getColumn(), c.table, where)
	return c.GetData(query, sid, date, num)
}

func (c *Chat) SetData() []byte {
	uid := c.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Chat.SetData'uid is nil", nil)
	}
	sid := c.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Chat.SetData'sid is nil", nil)
	}

	content := c.values.Get("content")
	if len(content) == 0 {
		return Result(0, "Chat.SetData'sid is nil", nil)
	}
	user := c.values.Get("user")
	if i, err := strconv.Atoi(user); err != nil || i == 0 {
		return Result(0, "Chat.SetData'user is 0", nil)
	}
	username := c.values.Get("username")
	analysis := c.values.Get("analysis")
	at := c.values.Get("at")

	tx, err := db.Begin()
	if err != nil {
		return Result(0, err.Error(), nil)
	}

	s := fmt.Sprintf("INSERT INTO jyh_record (`cid`, `user`, `status`) (SELECT '%s', `user`, 1 FROM jyh_mate WHERE `sid` = '%s') ON DUPLICATE KEY UPDATE `status`= VALUES(`status`)", uid, sid)
	// fmt.Println(s)
	db.Exec("INSERT INTO jyh_chat (`uid`, `sid`, `user`, `username` , `content`, `analysis`, `at`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, 1) ON DUPLICATE KEY UPDATE `sid` = VALUES(`sid`), `user`=VALUES(`user`), `username`=VALUES(`username`), `content`=VALUES(`content`), `analysis`=VALUES(`analysis`) ,`at`=VALUES(`at`)", uid, sid, user, username, content, analysis, at)
	db.Exec(s)
	db.Exec("DELETE FROM jyh_record WHERE `data` < DATE_ADD(NOW(), INTERVAL -1 YEAR)")
	err = tx.Commit()
	if err != nil {
		return Result(0, err.Error(), nil)
	}
	return Result(1, "chat insert is ok", nil)
}
