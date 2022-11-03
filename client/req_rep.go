package client

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/pebbe/zmq4"
)

type ReqRepClient struct {
	Client
	Timeout time.Duration
}

type ReqRepRequest struct {
	Endpoint string      `json:"type"`
	Data     interface{} `json:"data"`
}

type ReqRepResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func NewReqRepClient(host, port string, protocol Protocol, context *zmq4.Context, timeout time.Duration) (*ReqRepClient, error) {
	client, err := NewClient(host, port, protocol, zmq4.REQ, context)
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		client.Socket.SetRcvtimeo(timeout)
	}
	return &ReqRepClient{
		Client:  *client,
		Timeout: timeout,
	}, nil
}

func (c *ReqRepClient) Request(endpoint string, data interface{}) (*ReqRepResponse, error) {
	if !c.isConnected {
		return nil, errors.New("socket is not connected")
	}
	message, err := json.Marshal(ReqRepRequest{
		Endpoint: endpoint,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}
	_, err = c.Socket.SendBytes(message, 0)
	if err != nil {
		return nil, err
	}
	responseData, err := c.Socket.RecvBytes(0)
	if err != nil {
		return nil, err
	}
	response := ReqRepResponse{}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (c ReqRepClient) RequestRaw(endpoint string, data interface{}) ([]byte, error) {
	if !c.isConnected {
		return nil, errors.New("socket is not connected")
	}
	message, err := json.Marshal(ReqRepRequest{
		Endpoint: endpoint,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}
	_, err = c.Socket.SendBytes(message, 0)
	if err != nil {
		return nil, err
	}
	responseData, err := c.Socket.RecvBytes(0)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
