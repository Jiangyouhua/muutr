package model

import (
	"fmt"
	"net/url"
	"encoding/json"
	"errors"
	"strings"
	"strconv"
)

type Sync struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_Sync",
		columns: nil,
		primary: "aid",
	}
	s := &Sync{m}
	models["Sync"] = s
}

func GetSync(v url.Values) *Sync {
	m := models["Sync"].(*Sync)
	m.values = v
	return m
}

// 路由
func (s *Sync) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Sync.Routing'key arg is nil", nil)
	}
	s.values = v
	m := map[string]func() []byte{
		"GetBase":   s.GetBase,
		"GetData":   s.GetAll,
		"GetSingle": s.GetSingle,
		"SetData":   s.SetAll,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("sync is not this func (%s)", k), nil)
	}
	return f()
}

func (s *Sync) SetAndGet() (data map[string]interface{}, err error) {
	err = s.SetIts()
	if err != nil {
		return
	}
	return s.GetIts()
}

// 获取全部
func (s *Sync) GetIts() (data map[string]interface{}, err error) {
	user := Session.Get("user").(string)
	data = make(map[string]interface{})
	text := make([]string, 0) // 错误信息汇总

	// 是否为同步
	var editor string
	if t, err := strconv.Atoi(s.values.Get("type")); err == nil && t == 1 {
		editor = fmt.Sprintf("AND `editor` != %s", user)
	}

	// getJoin
	date := s.values.Get("join")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT * FROM jyh_join WHERE `user` = ? AND `date` > ? %s ", editor)
	d, err := s.Fetch(query, user, date)
	if err != nil {
		text = append(text, "Join is err: "+err.Error())
	} else {
		data["join"] = d
	}

	// getMoney
	date = s.values.Get("money")
	ymd, _ := strconv.Atoi(date)
	y := ymd / 10000
	query = GetNote(nil).QueryMoneyGroupByYearMonth(user, y)
	d, err = s.Fetch(query)
	if err != nil {
		text = append(text, "Money is err: "+err.Error())
	} else {
		data["money"] = d
	}

	// getRecord
	date = s.values.Get("recode")
	if len(date) == 0 {
		date = datetime
	}
	query = fmt.Sprintf("SELECT * FROM jyh_record WHERE `user` = ? AND `date` > ? %s", editor)
	d, err = s.Fetch(query, user, date)
	if err != nil {
		text = append(text, "recode is err: "+err.Error())
	} else {
		data["record"] = d
	}

	// get task
	date = s.values.Get("claim")
	if len(date) == 0 {
		date = datetime
	}
	query = fmt.Sprintf("SELECT * FROM jyh_claim WHERE `user` = ? AND `date` > ? %s", editor)
	d, err = s.Fetch(query, user, date)
	if err != nil {
		text = append(text, "Claim is err: "+err.Error())
	} else {
		data["claim"] = d
	}

	// 获取分类相关数据
	query = "SELECT `sid`, `uid` FROM jyh_join WHERE `user` = ? AND `status` = 1 "
	d, err = s.Fetch(query, user)
	if err != nil {
		text = append(text, "Get Sid and Uid is Err: "+err.Error())
		err = errors.New(strings.Join(text, ","))
		return
	}
	m := ReToString(d, "sid", "uid")
	if len(m["sid"]) == 0 {
		text = append(text, "Join's sid is nil")
		err = errors.New(strings.Join(text, ","))
		return
	}
	s.DataWithSorts(m["sid"], editor, data, text)
	return data, errors.New(strings.Join(text, ","))
}

// 根据用户加入的IDS返回需要同步到本地的数据
func (s *Sync) DataWithSorts(sids, editor string, data map[string]interface{}, text []string) {

	// get Sort related
	if len(sids) == 0 {
		str := s.values.Get("sids")
		sids = fmt.Sprintf("'%s'", str)
	}

	if len(sids) == 0 {
		text = append(text, "Sort sids or date is nil")
		return
	}

	//  get sort
	date := s.values.Get("sort")
	if len(date) == 0 {
		date = datetime
	}
	query := fmt.Sprintf("SELECT * FROM jyh_sort WHERE `uid` IN (%s) AND `date` > ? %s", sids, editor)
	d, err := s.Fetch(query, date)
	if err != nil {
		text = append(text, "Sort is err: "+err.Error())
	} else {
		data["sort"] = d
	}

	// get mate
	date = s.values.Get("mate")
	if len(date) == 0 {
		date = datetime
	}

	query = fmt.Sprintf("SELECT * FROM jyh_mate WHERE `sid` IN (%s) AND `date` > ? %s", sids, editor)
	d, err = s.Fetch(query, date)
	if err != nil {
		text = append(text, "Mate is err: "+err.Error())
	} else {
		data["mate"] = d
	}

	// get item
	date = s.values.Get("item")
	if len(date) == 0 {
		date = datetime
	}
	query = fmt.Sprintf("SELECT * FROM jyh_item WHERE `sid` IN (%s) AND `date` > ? %s", sids, editor)
	d, err = s.Fetch(query, date)
	if err != nil {
		text = append(text, "Item is err: "+err.Error())
	} else {
		data["item"] = d
	}

	// 获取任务项相关数据
	query = fmt.Sprintf("SELECT `uid` FROM jyh_item WHERE  `sid` IN (%s)", sids)
	d, err = s.Fetch(query)
	if err != nil {
		text = append(text, "ItemCondition is er: "+err.Error())
		err = errors.New(strings.Join(text, ","))
		return
	}

	m := ReToString(d, "uid")
	if len(m["uid"]) == 0 {
		text = append(text, "Item's uid is nil")
		err = errors.New(strings.Join(text, ","))
		return
	}

	// get node
	date = s.values.Get("node")
	if len(date) == 0 {
		date = datetime
	}
	query = fmt.Sprintf("SELECT * FROM jyh_node WHERE `iid` IN (%s) AND `date` > ? %s", m["uid"], editor)
	d, err = s.Fetch(query, date)
	if err != nil {
		text = append(text, "Node is err: "+err.Error())
	} else {
		data["node"] = d
	}
}

// 获取基本信息
func (s *Sync) GetBase() []byte {
	data := make(map[string]interface{})
	text := make([]string, 0) // 错误信息汇总

	// config
	date := s.values.Get("config")
	if len(date) == 0 {
		date = datetime
	}
	query := "SELECT * FROM jyh_config WHERE `uid` = 1 AND `date` > ?"
	d, err := s.Fetch(query, date)
	if err != nil {
		text = append(text, "Config is err: "+err.Error())
	} else {
		data["config"] = d
	}

	// ad
	date = s.values.Get("ad")
	if len(date) == 0 {
		date = datetime
	}
	query = "SELECT * FROM jyh_ad WHERE `date` > ?"
	d, err = s.Fetch(query, date)
	if err != nil {
		text = append(text, "Ad is err: "+err.Error())
	} else {
		data["ad"] = d
	}

	// holiday
	date = s.values.Get("holiday")
	if len(date) == 0 {
		date = datetime
	}
	query = "SELECT * FROM jyh_holiday WHERE `date` > ?"
	d, err = s.Fetch(query, date)
	if err != nil {
		text = append(text, "Holiday is err: "+err.Error())
	} else {
		data["holiday"] = d
	}
	return Result(1, strings.Join(text, ", "), data)
}

func (s *Sync) GetAll() []byte {
	data, err := s.GetIts()
	return Result(1, err.Error(), data)
}

func (s *Sync) GetSingle() []byte {
	data := make(map[string]interface{})
	text := make([]string, 0)
	s.DataWithSorts("", "", data, text)
	return Result(1, strings.Join(text, ",	"), data)
}

func (s *Sync) SetAll() []byte {
	err := s.SetIts()
	if err != nil {
		return Result(0, err.Error(), nil)
	}
	return Result(1, "SetAll is ok", nil)
}

func (s *Sync) SetIts() error {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return errors.New("Sync.SetAll'user  is nil")
	}

	var (
		data  interface{}
		err   string
		query string
	)

	d := s.values.Get("joins")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		join := GetJoin(s.values)
		query = join.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = e.Error()
		}
	}

	d = s.values.Get("sorts")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		sort := GetSort(s.values)
		query = sort.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("mates")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		mate := GetMate(s.values)
		query = mate.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("claims")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		task := GetClaim(s.values)
		query = task.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("items")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		item := GetItem(s.values)
		query = item.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("nodes")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		node := GetNode(s.values)
		query = node.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("records")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		record := GetRecord(s.values)
		query = record.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("notes")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		note := GetNote(s.values)
		query = note.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	d = s.values.Get("chats")
	if len(d) > 0 {
		json.Unmarshal([]byte(d), &data)
		chat := GetChat(s.values)
		query = chat.InsertUpdateStr(data)
		_, _, e := s.Set(query)
		if e != nil {
			err = fmt.Sprintf("%s, %s", err, e)
		}
	}

	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}
