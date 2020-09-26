package main

type room struct {

	// forwardは他のclientに送信するメッセージを保持するためのチャネル
	forward chan []byte

	join    chan *client     // チャットに参加しようとしているclientのためのチャネル
	leave   chan *client     // チャットから退出しようとしているclientのためのチャネル
	clients map[*client]bool // 在籍しているclientを保持
}

func (r *room) run() {
	for {
		select {
		// 参加
		case client := <-r.join:
			r.clients[client] = true
		// 退出
		case client := <-r.leave:
			delete(r.clients, client)
		// メッセージを送信
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}
