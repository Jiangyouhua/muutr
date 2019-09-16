package model

import (
	"net/url"
	"fmt"
	"strconv"
	"time"
)

type Note struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_note",
		columns: []string{"uid", "sid", "iid", "first", "date", "user", "username", "content", "money", "analysis", "editor", "status"},
		primary: "uid",
	}
	n := &Note{m}
	models["Note"] = n
}

func GetNote(v url.Values) *Note {
	m := models["Note"].(*Note)
	m.values = v
	return m
}

// 路由
func (n *Note) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Note.Routing'key arg is nil", nil)
	}
	n.values = v
	m := map[string]func() []byte{
		"GetWithUid":                        n.GetWithUid,
		"GetWithIidAndUserStartDate":        n.GetWithIidAndUserStartDate,
		"GetWithIidAndUserEndDate":          n.GetWithIidAndUserEndDate,
		"GetDisableWithSidAndUserStartDate": n.GetDisableWithSidAndUserStartDate,
		"GetDisableWithSidAndUserEndDate":   n.GetDisableWithSidAndUserEndDate,
		"GetMoneyGroupByMonth":              n.GetMoneyGroupByYearMonth,
		"SetData":                           n.SetDataWithPrimary,
		"GetData":                           n.GetDataWithPrimary,
		"DelData":                           n.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Note is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单条笔记信息
func (n *Note) GetWithUid() [] byte {
	uid := n.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Note.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", n.getColumn(), n.table)
	return n.GetData(query, uid)
}

// 获取用户目标笔记，最新N条信息
func (n *Note) GetWithIidAndUserStartDate() [] byte {
	iid := n.values.Get("iid")
	if len(iid) == 0 {
		return Result(0, "Note.GetWithIidAndUserStartDate'iid is nil", nil)
	}

	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Note.GetWithIidAndUserStartDate'user is nil", nil)
	}

	num := n.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "100"
	}

	date := n.values.Get("date")
	if len(date) == 0 {
		date = datetime + " 000"
	} else {
		num = "1000"
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE `iid` = ? AND `user` = ? AND `date` > ? ORDER BY `date` DESC LIMIT ? ", n.getColumn(), n.table)
	return n.GetData(query, iid, user, date, num)
}

// 分页，以NID为起点，获取群组内最新N条聊天数据
func (n *Note) GetWithIidAndUserEndDate() []byte {
	date := n.values.Get("date")
	if len(date) == 0 {
		return Result(0, "Note.GetWithIidAndUserEndDate'date is nil", nil)
	}
	iid := n.values.Get("iid")
	if len(iid) == 0 {
		return Result(0, "Note.GetWithIidAndUserEndDate'iid is nil", nil)
	}
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Note.GetWithIidAndUserEndDate'user is nil", nil)
	}
	num := n.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "200"
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `iid` = ? AND `user` = ? AND `date` < ? AND `status` = 1 ORDER BY `first` DESC LIMIT ? ", n.getColumn(), n.table);
	return n.GetData(query, iid, user, date, num)
}

// 获取用户目标笔记，最新N条信息
func (n *Note) GetDisableWithSidAndUserStartDate() [] byte {
	sid := n.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Note.GetWithSidAndUserStartDate'iid is nil", nil)
	}

	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Note.GetWithSidAndUserStartDate'user is nil", nil)
	}

	num := n.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "100"
	}

	date := n.values.Get("date")
	if len(date) == 0 {
		date = datetime + " 000"
	} else {
		num = "1000"
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE `status` = 0 AND `iid` in ( SELECT iid FROM jyh_item WHERE sid = ? AND `status` = 0) AND `user` = ? AND `date` > ? ORDER BY `date` DESC LIMIT ? ", n.getColumn(), n.table)
	return n.GetData(query, sid, user, date, num)
}

// 分页，以NID为起点，获取群组内最新N条聊天数据
func (n *Note) GetDisableWithSidAndUserEndDate() []byte {
	date := n.values.Get("date")
	if len(date) == 0 {
		return Result(0, "Note.GetWithSidAndUserEndDate'date is nil", nil)
	}
	sid := n.values.Get("sid")
	if len(sid) == 0 {
		return Result(0, "Note.GetWithSidAndUserEndDate'iid is nil", nil)
	}
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Note.GetWithSidAndUserEndDate'user is nil", nil)
	}
	num := n.values.Get("num")
	if i, err := strconv.Atoi(num); err != nil || i < 10 {
		num = "200"
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `status` = 0 AND `iid` in ( SELECT iid FROM jyh_item WHERE sid = ? AND `status` = 0) AND `user` = ? AND `date` < ? AND `status` = 1 ORDER BY `first` DESC LIMIT ? ", n.getColumn(), n.table);
	return n.GetData(query, sid, user, date, num)
}

// 请求一次
// 获取用户各月费用, 获取今年或去年的，或某一年的
func (n *Note) GetMoneyGroupByYearMonth() [] byte {
	year, _ := strconv.Atoi(n.values.Get("year"))
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Note.GetMoneyGroupByMonth'user is nil", nil)
	}
	query :=  n.QueryMoneyGroupByYearMonth(user, year)
	return n.GetData(query)
}

func(n *Note)QueryMoneyGroupByYearMonth(user string, year int)string{
	old := year
	if year < 2018 {
		year = time.Now().Year()
	}

	if year == time.Now().Year() {
		old = year - 1
	}
	return fmt.Sprintf(  "SELECT * FROM (" +
		" SELECT `sid`, `iid`, " +
		" SUM(CASE WHEN `money` > 0 THEN `money`*100 ELSE 0 END) 'income'," +
		" SUM(CASE WHEN `money` < 0 THEN `money`*100 ELSE 0 END) 'expenses'," +
		" SUM(`money`) 'surplus', EXTRACT(YEAR_MONTH FROM `first`) date" +
		" FROM jyh_note WHERE `status` = 1 AND (YEAR(`first`) = %v || YEAR(`first`) = %v) AND user = '%s'" +
		" GROUP BY EXTRACT(YEAR_MONTH FROM `first`), sid, iid" +
		") m ORDER BY date DESC",year, old, user)

}
