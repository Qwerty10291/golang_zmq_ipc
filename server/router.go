package server

import (zmq "github.com/pebbe/zmq4")

type RouterServer struct{
	Server
}

func NewRouterServer(host, port string, protocol Protocol, context *zmq.Context) (*RouterServer, error) {
	server, err := NewServer(host, port, protocol, zmq.ROUTER, context)
	if err != nil{
		return nil, err
	}
	return &RouterServer{Server: *server}, nil
}
