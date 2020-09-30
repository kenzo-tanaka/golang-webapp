package main

import (
	"golang-webapp/chap1/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {

	// forwardは他のclientに送信するメッセージを保持するためのチャネル
	forward chan *message

	join    chan *client     // チャットに参加しようとしているclientのためのチャネル
	leave   chan *client     // チャットから退出しようとしているclientのためのチャネル
	clients map[*client]bool // 在籍しているclientを保持
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		// 参加
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("私いクライアントが参加しました。")
		// 退出
		case client := <-r.leave:
			delete(r.clients, client)
			r.tracer.Trace("クライアントが退出しました。")
		// メッセージを送信
		case msg := <-r.forward:
			r.tracer.Trace("メッセージが送信されました。", string(msg.Message))
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()

	// goroutineを使って別のスレッドを割当
	go client.write()

	// メインのスレッドではreadメソッドを呼び出ししている
	client.read()
}
