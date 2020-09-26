package main

type room struct {

	// forwardは他のclientに送信するメッセージを保持するためのチャネル
	forward chan []byte

	join    chan *client     // チャットに参加しようとしているclientのためのチャネル
	leave   chan *client     // チャットから退出しようとしているclientのためのチャネル
	clients map[*client]bool // 在籍しているclientを保持
}
