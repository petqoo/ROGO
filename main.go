package main

import (
	// "fmt"
	"log/slog"
	"net"
	"strconv"
)


type Config struct{
	Port int 
}

type Server struct{
	Config Config
	Peers map[*Peer]bool
	ln net.Listener
	addpeerch chan *Peer
	stopchannel chan struct{}
	msgchannel chan string
    cmdchann chan Command
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
		stopchannel: make(chan struct{}),
		msgchannel: make(chan string, 100),
		cmdchann: make(chan Command,100),
	}
}
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgchannel:
			slog.Info("Received message", "message", msg)
		case <-s.stopchannel:
			slog.Info("Stopping server...")
			// for peer := range s.Peers {
			// 	slog.Info("Closing peer connection", "peer", peer)
				
			// }
			return 
		case peer := <-s.addpeerch:
			s.Peers[peer] = true
			slog.Info("New peer added", "peer", peer)

        case cmd := <-s.cmdchann:
			slog.Info("Received command", "command", cmd)
			switch cmd := cmd.(type) {
			case SetCommand:
				slog.Info("Processing SET command", "key", cmd.Key, "value", cmd.Value)
			}
		}
	}
}
func (s *Server) Start() error {
	ln,err:=net.Listen("tcp", ":"+strconv.Itoa(s.Config.Port))
	if err !=nil {
		slog.Info("Failed to start server", "error", err)
		return  err
	}

	s.ln= ln

	go s.loop()
	slog.Info("Server started", "port", s.Config.Port)
	return s.acceptloop()
	
}
func (s *Server) acceptloop() error {
	for {
		conn,err:= s.ln.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}
		slog.Info("Accepted new connection", "remote_addr", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn){
 peers := NewPeer(&conn, s.msgchannel,s.cmdchann)
 s.addpeerch <- peers
 slog.Debug("raw ja jdid")
 peers.Reedloop()

}
func main (){
	slog.Info("Starting server...")
	server:= NewServer(NewConfig(8085))
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}


