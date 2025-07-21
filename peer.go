package main

import "net"

type Peer struct {
	conn *net.Conn
}

func NewPeer(conn *net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}
func(p *Peer) Reedloop(){

}