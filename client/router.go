package client

import (
	"time"
	"github.com/pebbe/zmq4"
)

type RouterClient struct {
	ReqRepClient
}

func NewRouterClient(host, port string, protocol Protocol, context *zmq4.Context, timeout time.Duration) (*RouterClient, error) {
	client, err := NewClient(host, port, protocol, zmq4.REQ, context)
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		client.Socket.SetRcvtimeo(timeout)
	}

	return &RouterClient{
		ReqRepClient{
			Client:  *client,
			Timeout: timeout,
		}}, nil
}
