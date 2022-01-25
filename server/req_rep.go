package server

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

type ReqRepServer struct{
	Server
	handlers map[string]ReqRepServerHandler
}

type ReqRepResponse struct{
	Status bool `json:"status"`
	Data interface{} `json:"data"`
	Error string `json:"error"`
}

type reqRepRequest struct{
	Endpoint string      `json:"type"`
	Data     interface{} `json:"data"`
}

type ReqRepServerHandler func(interface{}) ReqRepResponse


func NewReqRepServer(host, port string, protocol Protocol, context *zmq.Context) (*ReqRepServer, error){
	server, err := NewServer(host, port, protocol, zmq.REP, context)
	if err != nil{
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
	if !s.isBinded{
		return errors.New("socket is not binded")
	}
	go s.Listener()
	return nil
}

func (s *ReqRepServer) Listener() {
	for s.isBinded{
		data, err := s.socket.RecvBytes(0)
		if err != nil{
			panic(err)
		}
		request :=  reqRepRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil{
			panic(err)
		}
		if handler, ok := s.handlers[request.Endpoint]; ok{
			resp, err := json.Marshal(handler(request))
			if err != nil{
				panic(err)
			}
			_, err = s.socket.SendBytes(resp, 0)
			if err != nil{
				panic(err)
			}
		}
	}
}