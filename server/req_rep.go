package server

import (
	"encoding/json"
	"errors"
	"log"

	zmq "github.com/pebbe/zmq4"
)

type ReqRepServer struct {
	Server
	handlers map[string]ReqRepServerHandler
}

type ReqRepResponse interface{}

type reqRepRequest struct {
	Endpoint string      `json:"type"`
	Data     interface{} `json:"data"`
}

type ReqRepServerHandler func(interface{}) ReqRepResponse

func NewReqRepServer(host, port string, protocol Protocol, context *zmq.Context) (*ReqRepServer, error) {
	server, err := NewServer(host, port, protocol, zmq.REP, context)
	if err != nil {
		return nil, err
	}
	return &ReqRepServer{
		Server:   *server,
		handlers: map[string]ReqRepServerHandler{},
	}, nil
}

func (s *ReqRepServer) NewHandler(endpoint string, handler ReqRepServerHandler) {
	s.handlers[endpoint] = handler
}

func (s *ReqRepServer) Start() error {
	if !s.isBinded {
		return errors.New("socket is not binded")
	}
	go s.Listener()
	return nil
}

func (s *ReqRepServer) Listener() {
	for s.isBinded {
		data, err := s.Socket.RecvBytes(0)
		if err != nil {
			log.Printf("failed to recv bytes from %s:%s", s.Port, err)
			continue
		}

		request := reqRepRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			log.Printf("failed to parse json request on socket %s:%s", s.Port, err)
			continue
		}

		if handler, ok := s.handlers[request.Endpoint]; ok {
			resp, err := json.Marshal(handler(request))
			if err != nil {
				panic(err)
			}
			_, err = s.Socket.SendBytes(resp, 0)
			if err != nil {
				log.Printf("failed to send response from %s:%s", s.Port, err)
				continue
			}
		}
		
	}
}
