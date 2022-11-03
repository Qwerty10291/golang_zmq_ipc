package client

import (
	zmq "github.com/pebbe/zmq4"
)

type DealerClient struct {
	Client
}

func NewDealerClient(host, port string, protocol Protocol, contex *zmq.Context) (*DealerClient, error) {
	client, err := NewClient(host, port, protocol, zmq.REP, contex)
	if err != nil {
		return nil, err
	}
	return &DealerClient{
		Client:        *client,
	}, nil
}