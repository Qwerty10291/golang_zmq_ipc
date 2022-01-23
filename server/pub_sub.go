package server

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

type PubSubServer struct{
	Server
}

type pubSubFrame struct {
	Topic string `json:"topic"`
	Data interface{} `json:"data"`
}

func NewPubSubServer(host, port string, protocol Protocol, context *zmq.Context,) (*PubSubServer, error) {
	server, err := NewServer(host, port, protocol, zmq.PUB, context)
	if err != nil{
		return nil, err
	}
	return &PubSubServer{
		Server: *server,
	}, nil
}


func (s *PubSubServer) Send(topic string, message interface{}) error{
	if !s.isBinded{
		return errors.New("socket is not binded")
	}
	frame := pubSubFrame{
		Topic: topic,
		Data:  message,
	}

	data, err := json.Marshal(frame)
	if err != nil{
		return err
	}
	_, err = s.socket.SendBytes(data, 0)
	return err
}
