// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"encoding/json"
	"ido/logger"
	"strconv"
	"fmt"
)

const (
	// Time allowed to write a Message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong Message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum Message size allowed from peer.
	maxMessageSize = 512
)

var (
	Log *logger.Logger
	newline = []byte{'\n'}
	space   = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Client is a middleman between the webSocket connection and the hub.
type Client struct {
	hub *Hub

	// The webSocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// 来自哪个聊天室
	sid string

	// 用户ID
	user string
}

// readPump pumps messages from the webSocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// 读取数据内容
		_, me, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				Log.Out(logger.RUN_LOG, fmt.Sprintf("ReadMessage: %s, %s", err.Error(), me))
			}
			break
		}
		//Log.Out(logger.RUN_LOG, fmt.Sprintf("SocketRead: %s", me))
		// 解析字符串
		// 聊天信息
		m := new(Message)
		if err := json.Unmarshal(me, m); err != nil || m == nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				Log.Out(logger.ERR_LOG, fmt.Sprintf("Json>Unmarshal: %s, %s", err.Error(), m))

			}
			break
		}
		m.Date =strconv.FormatInt(time.Now().UnixNano()/1000,10)
		SaveToChat(m)
		c.hub.message <- m
	}
}

// writePump pumps messages from the hub to the webSocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current webSocket Message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles webSocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return
	}
	sid := r.Form.Get("sid")
	user := r.Form.Get("user")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), sid: sid, user: user}
	Log.Out(logger.RUN_LOG, fmt.Sprintf("ServWs: %s, %s", sid, user))
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
