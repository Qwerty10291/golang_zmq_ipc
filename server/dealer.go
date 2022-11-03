package server

import (zmq "github.com/pebbe/zmq4")

type DealerServer struct{
	Server
}

func NewDealerServer(host, port string, protocol Protocol, context *zmq.Context) (*DealerServer, error) {
	server, err := NewServer(host, port, protocol, zmq.DEALER, context)
	if err != nil{
		return nil, err
	}
	return &DealerServer{Server: *server}, nil	
}
