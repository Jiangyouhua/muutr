package model

import (
	"net/url"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type Holiday struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_holiday",
		columns: []string{"uid", "year", "month", "day", "date", "status"},
		primary: "uid",
	}
	h := &Holiday{m}
	models["Holiday"] = h
}

func GetHoliday(v url.Values)*Holiday{
	m := models["Holiday"].(*Holiday)
	m.values = v
	return m
}

// 路由
func (h *Holiday) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Holiday.Routing'key arg is nil", nil)
	}
	h.values = v
	m := map[string]func() []byte{
		"GetWithUid":  h.GetWithUid,
		"GetWithYear": h.GetWithYear,
		"SetWithUid":  h.SetWithUid,
		"SetData":     h.SetDataWithPrimary,
		"GetData":     h.GetDataWithPrimary,
		"DelData":     h.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Holiday is not this func (%s)", k), nil)
	}
	return f()
}

// 获取某天的放假信息
func (h *Holiday) GetWithUid() []byte {
	uid := h.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Holiday.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", h.getColumn(), h.table)
	return h.GetData(query, uid)
}

// 获取某年的放假信息
func (h *Holiday) GetWithYear() []byte {
	year := h.values.Get("year")
	if len(year) == 0 {
		t := time.Now();
		year = strconv.Itoa(t.Year())
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `year` = ?", h.getColumn(), h.table)
	return h.GetData(query, year)
}

func (h *Holiday) SetWithUid() []byte {
	uid := h.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Holiday.SetWithUid'uid is nil", nil)
	}
	arr := strings.Split(uid, "-")
	if len(arr) != 3 {
		return Result(0, "Holiday.SetWithUid'uid is err", nil)
	}
	h.values.Set("year", arr[0])
	h.values.Set("month", arr[1])
	h.values.Set("day", arr[2])
	h.values.Set("uid", strings.Join(arr, ""))
	return h.SetDataWithPrimary();
}
