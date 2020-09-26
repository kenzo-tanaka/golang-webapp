package main

type room struct {

	// forwardは他のclientに送信するメッセージを保持するためのチャネル
	forward chan []byte

	join    chan *client
	leave   chan *client
	clients map[*client]bool
}
