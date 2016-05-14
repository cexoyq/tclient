// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	pub "t_client/pub"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	//"encoding/json"
	"github.com/gorilla/websocket"
	//linuxproc "github.com/c9s/goprocinfo/linux"
)

var addr = flag.String("addr", "192.168.1.112:80", "http service address")

func main() {
	var (
		c *websocket.Conn
	)
	
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		var (
			msg pub.Msg1001
			head pub.MsgHead
			body []pub.StreamInfo
		)
		for {
			//msg=pub.Msg1000{}
			//_, message, err := c.ReadMessage()
			//log.Printf("read:%s",message)
			err := c.ReadJSON(&msg)
			head = msg.Head
			body = msg.Body
			if err != nil {
				log.Println("read:", err)
				return
			}
			if(head.MsgType == 2){
				log.Printf("msg_r.MsgType == 2!")
				log.Printf("recv,ServerNum: %d;,MsgType: %d;,MsgLen: %d;",head.ServerNum, head.MsgType,head.MsgLen)
				log.Printf("recv,ServerNum: Body.Name:%+v;",body)
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()

	for {
			var (
				msg pub.Msg1001
				head pub.MsgHead
				body []pub.StreamInfo
			)
			body =[]pub.StreamInfo{
				{1,"1号桌_全景","rtsp://","rtmp://",},
				{2,"2号桌_全景","rtsp://2","rtmp://2",},
			}
			head = pub.MsgHead{
                1,
				2,
				3,
            }
			msg = pub.Msg1001{
				head,
				body,
			}
		select {
		case t := <-ticker.C:
			log.Printf("write:",msg)
			err := c.WriteJSON(msg)
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err,t.String())
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}