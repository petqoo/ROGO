package main

import (
	// "fmt"
	"log/slog"
	"net"
	"strconv"
	"sync"
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
	mu   sync.RWMutex
	data map[string][]byte
	delpeerch chan *Peer
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
		data: make(map[string][]byte),
		delpeerch: make(chan *Peer, 100),
	}
}
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgchannel:
			slog.Info("Received message", "message", msg)
		case <-s.stopchannel:
			s.mu.RLock()
			for peer := range s.Peers {
				peer.delpeerch <- peer
				slog.Info("Closing peer connection", "peer", peer)
			}
			s.mu.RUnlock()
			slog.Info("Stopping server...")
			return

		case peer := <-s.addpeerch:
			s.mu.Lock()
			s.Peers[peer] = true
			s.mu.Unlock()
			slog.Info("New peer added", "peer", peer)
        case peer := <-s.delpeerch:
			s.mu.Lock()
			if _, exists := s.Peers[peer]; exists {
				slog.Info("Removing peer", "peer", peer)
				delete(s.Peers, peer)
				(*peer.conn).Close()
				slog.Info("Peer connection closed", "peer", peer)
			} else {
				slog.Info("Peer not found for removal", "peer", peer)
			}
			s.mu.Unlock()
        case cmd := <-s.cmdchann:
			slog.Info("Received command", "command", cmd)
			switch cmd := cmd.(type) {
			case *SetCommand:
				slog.Info("Processing SET command", "key", cmd.Key, "value", cmd.Value)
				s.handleSetCommand(cmd)
			case *Getcommand:
				slog.Info("Processing GET command", "key", cmd.Key)
				s.mu.RLock()
				value, exists := s.data[cmd.Key]
				s.mu.RUnlock()
				if exists {
					(*cmd.Peer.conn).Write([]byte("VALUE " + cmd.Key + " " + string(value) + "\n"))
					slog.Info("Value retrieved", "key", cmd.Key, "value", string(value))
				} else {
					(*cmd.Peer.conn).Write([]byte("NOT_FOUND\n"))
					slog.Info("Key not found", "key", cmd.Key)
				}
			}

		}
	}
}
func (s *Server) handleSetCommand(cmd *SetCommand) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[cmd.Key] = []byte(cmd.Value)
	(*cmd.Peer.conn).Write([]byte("OK\n"))
	slog.Info("Data stored", "key", cmd.Key)
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
 peers := NewPeer(&conn, s.msgchannel,s.cmdchann, s.delpeerch)
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


