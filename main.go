package main

import (
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
}
func NewConfig(port int)Config{
	return  Config{
		Port: port,
	}
}

func NewServer(config Config ) *Server{
	return &Server{
		Config: config,
		Peers: make([]*Peer, 0),
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

}
func main (){
	slog.Info("Starting server...")
}

