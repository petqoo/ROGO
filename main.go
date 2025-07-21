package main

import (
	"fmt"
	"log/slog"
	"net"
)


type Config struct{
	Port int 
}

type Server struct{
	Config Config
	Peers map[*Peer]bool
	ln net.Listener
	addpeerch chan *Peer
}
func NewConfig(port int)Config{
	return  Config{
		Port: port,
	}
}

func NewServer(config Config ) *Server{
	return &Server{
		Config: config,
		Peers: make(map[*Peer]bool),
		addpeerch: make(chan *Peer, 100), 
	}
}
func (s *Server) loop() error{
	for {
		select {
		case peer := <-s.addpeerch:
			s.Peers[peer] = true
			slog.Info("New peer added", "peer", peer)
	    default:
			fmt.Println(" mjnoun nta za3ma")

		}
	}
}
func (s *Server) Start() error {
	ln,err:=net.Listen("tcp", ":"+string(s.Config.Port))
	if err !=nil {
		slog.Info("Failed to start server", "error", err)
		return  err
	}
	s.ln= ln
	return s.acceptloop()
	
}
func (s *Server) acceptloop() error {
	for {
		conn,err:= s.ln.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn){
 peers := NewPeer(&conn)
 s.addpeerch <- peers
 slog.Debug("raw ja jdid")
 peers.Reedloop()

}
func main (){
	slog.Info("Starting server...")
	server:= NewServer(NewConfig(8080))
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}


