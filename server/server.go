package server

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type Server struct {
	Host     string
	Port     string
	Protocol Protocol
	isBinded bool
	Socket   *zmq.Socket
}

func NewServer(host, port string, protocol Protocol, socketType zmq.Type, context *zmq.Context) (*Server, error) {
	socket, err := context.NewSocket(socketType)
	if err != nil {
		return nil, err
	}
	return &Server{
		Host:     host,
		Port:     port,
		Protocol: protocol,
		isBinded: false,
		Socket:   socket,
	}, nil
}

func (s *Server) Bind() error {
	err := s.Socket.Bind(fmt.Sprintf("%s://%s:%s", s.Protocol, s.Host, s.Port))
	if err != nil {
		return err
	}
	s.isBinded = true
	return nil
}

func (s *Server) Close() error {
	err := s.Socket.Close()
	if err != nil {
		return err
	}
	s.isBinded = false
	return nil
}


func (s *Server) GetSocket() *zmq.Socket {
	return s.Socket
}