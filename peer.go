package main

import (
	"io"
	"net"
)

type Peer struct {
	conn *net.Conn
	msgch chan string
	cmdcha chan  Command
	delpeerch chan *Peer
}

func NewPeer(conn *net.Conn, msgch chan string, cmdch chan Command,delpeerch chan *Peer) *Peer {
	return &Peer{
		conn: conn,
		msgch: msgch,
		cmdcha: cmdch,
		delpeerch: delpeerch,
	}
}
func(p *Peer) Reedloop(){
	for {
		buf := make([]byte, 1024)
		n, err := (*p.conn).Read(buf)
		if err != nil {
			if err==io.EOF{
				println("Connection closed by peer")				
			}
			p.delpeerch <- p
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
			if cmd == nil {
				println("ml79ni wal0")
				continue
			}
			//type assertion khra
			switch cmd.(type) {
			case *Getcommand:
				println("Received GET command")
				cmd.(*Getcommand).Peer = p
			case *SetCommand:
				println("Received SET command")
				cmd.(*SetCommand).Peer=p
			default:
				println("Unknown command type")
			}
			
			p.cmdcha <-cmd
			
			println("Parsed command:", cmd)
		}
	}

}