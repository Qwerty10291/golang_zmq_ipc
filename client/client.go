package client

import (
	"errors"
	"fmt"

	"github.com/pebbe/zmq4"
)

type Client struct {
	Host     string
	Port     string
	Protocol Protocol
	isConnected bool

	socket *zmq4.Socket
}

func NewClient(host, port string, protocol Protocol, socketType zmq4.Type, contex *zmq4.Context) (*Client, error){
	socket, err := contex.NewSocket(socketType)
	if err != nil {
		return nil, err
	}
	SetSocketFlags(socket)
	return &Client{
		Host:        host,
		Port:        port,
		Protocol:    protocol,
		isConnected: false,
		socket:      socket,
	}, nil
}

func (c *Client) Connect() error{
	if c.isConnected{
		return errors.New("socket already connected")
	}
	err := c.socket.Connect(fmt.Sprintf("%s://%s:%s", c.Protocol, c.Host, c.Port))
	if err != nil{
		return err
	}
	c.isConnected = true
	return nil
}

func (c *Client) Close() error{
	if !c.isConnected{
		return errors.New("socket is not connected")
	}
	err := c.socket.Close()
	if err != nil{
		return err
	}
	c.isConnected = false
	return nil
}