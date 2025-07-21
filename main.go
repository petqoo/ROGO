package main

import "net"


type Config struct{
	Port int 
}

type Server struct{
	Config Config
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
	}
}
func (s *Server) Start() error {
	ln,err:=net.Listen("tcp", ":"+string(s.Config.Port))
	if err !=nil {
		return  err
	}
	s.ln= ln
	return nil
	
}
func (s *Server) Runloop(
	hey,err:= 
)