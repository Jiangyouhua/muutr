// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	//"strings"
	"net/url"
	"ido/model"
	"ido/logger"
	"strings"
)

// hub maintains the set of active clients and messages messages to the
// clients.
// 接收的信息类型
type Message struct {
	Uid      string `json:"uid"`
	Sid      string `json:"sid"`
	Date     string `json:"date"`
	User    string `json:"user"`
	Content  string `json:"content"`
	Username string `json:"username"`
	At       string `json:"at"`
	Analysis string `json:"analysis"`
	Status   string `json:"status"`
}

// 分聊天室接收连接，手机号对应各客户端
type clients map[string]*Client

// webSocket 接收主休
type Hub struct {
	// 客户端, 映射到各群组（各聊天室）.
	clients map[string]clients

	// 聊天的信息
	message chan *Message

	// 注册客户端
	register chan *Client

	// 注销客户端
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		message:      make(chan *Message),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[string]clients),
	}
}

// 保存对话到数据库
func SaveToChat(m *Message) {
	v := url.Values{}
	v.Set("handle", "Chat")
	v.Set("func", "SetData")
	v.Set("uid", m.Uid)
	v.Set("sid", m.Sid)
	v.Set("user", m.User)
	v.Set("username", m.Username)
	v.Set("content", m.Content)
	v.Set("analysis", m.Analysis)
	v.Set("status", m.Status)
	v.Set("date", m.Date)
	v.Set("at", m.At)
	model.Router(v)
}

// 注册、注销用户连接，接收聊天信息
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			sid := client.sid
			c, ok := h.clients[sid]
			if !ok {
				c = make(map[string]*Client)
				c[client.user] = client
				h.clients[sid] = c
			} else {
				c[client.user] = client
			}
		case client := <-h.unregister:
			sid := client.sid
			c, ok := h.clients[sid]
			if !ok {
				break
			}
			if _, ok := c[client.user]; ok {
				delete(c, client.user)
				close(client.send)
			}
			if len(c) == 0 {
				delete(h.clients, sid)
			}
		case message := <-h.message: // 发送信息
			// 发送信息到各用户
			sid := message.Sid
			c, ok := h.clients[sid]
			if !ok {
				Log.Out(logger.RUN_LOG, "hub.Run(case message) sid is nil")
				break
			}
			b, err := json.Marshal(message)
			if err != nil {
				Log.Out(logger.ERR_LOG, err.Error())
				b = []byte(err.Error())
			}

			// 发给所有人
			if len(message.At) == 0 {
				for _, client := range c {
					h.Send(client, b)
				}
				break
			}

			// 发给at的人, 加上自己
			arr := strings.Split(message.At, ",")
			arr = append(arr, message.User)
			for _, u := range arr {
				if client, ok := c[u]; ok && client != nil {
					h.Send(client, b)
				}
			}
		}
	}
}

func (h *Hub) Send(client *Client, b []byte){
	select {
	case client.send <- b:
	default:
		close(client.send)
	}
}