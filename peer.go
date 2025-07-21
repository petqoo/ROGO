package main

import "net"

type Peer struct {
	conn *net.Conn
	msgch chan string

}

func NewPeer(conn *net.Conn, msgch chan string) *Peer {
	return &Peer{
		conn: conn,
		msgch: msgch,
	}
}
func(p *Peer) Reedloop(){
	for {
		buf := make([]byte, 1024)
		n, err := (*p.conn).Read(buf)
		if err != nil {
			return
		}
		if n > 0 {
			data := buf[:n]
			println("Received data:", string(data))
			p.msgch <- string(data)
		}
	}

}