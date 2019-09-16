package model

import (
	"net/url"
	"fmt"
)

type Node struct {
	Model
}

func init() {
	initModels()
	m := Model{
		table:   "jyh_node",
		columns: []string{"uid", "tid", "iid", "ymd", "user", "username", "first", "date", "editor", "status"},
		primary: "uid",
	}
	n := &Node{m}
	models["Node"] = n
}

func GetNode(v url.Values)*Node{
	m := models["Node"].(*Node)
	m.values = v
	return m
}

// 路由
func (n *Node) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "Node.Routing'key arg is nil", nil)
	}
	n.values = v
	m := map[string]func() []byte{
		"GetWithUid":  n.GetWithUid,
		"GetWithTid":  n.GetWithTid,
		"GetWithIid":  n.GetWithIid,
		"GetWithUser": n.GetWithUser,
		"SetData":     n.SetDataWithPrimary,
		"GetData":     n.GetDataWithPrimary,
		"DelData":     n.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("Node is not this func (%s)", k), nil)
	}
	return f()
}

// 获取某目标节点的单个用户完成信息
func (n *Node) GetWithUid() [] byte {
	uid := n.values.Get("uid")
	if len(uid) == 0 {
		return Result(0, "Node.GetWithUid'uid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `uid` = ?", n.getColumn(), n.table)
	return n.GetData(query, uid)
}

// 获取某目标节点的用户完成信息
func (n *Node) GetWithTid() [] byte {
	tid := n.values.Get("tid")
	if len(tid) == 0 {
		iid := n.values.Get("iid")
		if len(iid) == 0 {
			return Result(0, "Node.GetWithTid'uid or iid is nil", nil)
		}
		date := n.values.Get("date")
		if len(date) == 0 {
			return Result(0, "Node.GetWithTid'uid or date is nil", nil)
		}
		tid = fmt.Sprintf("%s,%s", iid, date)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `tid` = ?", n.getColumn(), n.table)
	return n.GetData(query, tid)
}

// 获取某目标的用户完成信息
func (n *Node) GetWithIid() [] byte {
	iid := n.values.Get("iid")
	if len(iid) == 0 {
		return Result(0, "Node.GetWithIid'iid is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `iid` = ?", n.getColumn(), n.table)
	return n.GetData(query, iid)
}

// 获取某目标的用户完成信息
func (n *Node) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "Node.GetWithIid'user is nil", nil)
	}
	//query := fmt.Sprintf("SELECT %s FROM %s WHERE `user` = ?", n.getColumn(), n.table)
	query := fmt.Sprintf("SELECT `username` FROM %s WHERE AND `user` = ?", n.table)
	return n.GetData(query, user)
}

