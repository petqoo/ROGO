package main

import "net"

type Peer struct {
	conn *net.Conn
	msgch chan string
	cmdcha chan  SetCommand
}

func NewPeer(conn *net.Conn, msgch chan string) *Peer {
	return &Peer{
		conn: conn,
		msgch: msgch,
		cmdcha: make(chan SetCommand, 100),
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
			cmd,err:= parseCommand(string(data))
			if err != nil {
				println("Error parsing command:", err.Error())
				continue
			}
			p.cmdcha <-cmd.( SetCommand)
			println("Parsed command:", cmd)
		}
	}

}