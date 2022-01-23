package server

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type Server struct {
	Host string
	Port string
	Protocol Protocol
	isBinded bool
	socket *zmq.Socket
}

func NewServer(host, port string, protocol Protocol, socketType zmq.Type, context *zmq.Context) (*Server, error) {
	socket, err := context.NewSocket(socketType)
	if err != nil{
		return nil, err
	}
	return &Server{
		Host:     host,
		Port:     port,
		Protocol: protocol,
		isBinded: false,
		socket:   socket,
	}, nil
}

func (s *Server) Bind() error {
	err := s.socket.Bind(fmt.Sprintf("%s://%s:%s", s.Protocol, s.Host, s.Port))
	if err != nil{
		return err
	}
	s.isBinded = true
	return nil
}

func (s *Server) Close() error{
	err := s.socket.Close()
	if err != nil{
		return err
	}
	s.isBinded = false
	return nil
}
