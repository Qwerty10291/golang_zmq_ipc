package client

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

type PubSubClient struct {
	Client
	topicHandlers map[string]PubSubMessageHandler
}

type pubSubFrame struct {
	Topic string
	Data  interface{}
}

type PubSubMessageHandler func(interface{})

func NewPubSubClient(host, port string, protocol Protocol, contex *zmq.Context) (*PubSubClient, error) {
	client, err := NewClient(host, port, protocol, zmq.SUB, contex)
	if err != nil {
		return nil, err
	}
	return &PubSubClient{
		Client:        *client,
		topicHandlers: map[string]PubSubMessageHandler{},
	}, nil
}

func (c *PubSubClient) AddTopicHandler(topic string, handler PubSubMessageHandler) {
	c.topicHandlers[topic] = handler
}

func (c *PubSubClient) Connect() error {
	err := c.Client.Connect()
	if err != nil {
		return err
	}
	err = c.Socket.SetSubscribe("")
	return err
}

func (c *PubSubClient) Start() error {
	if !c.isConnected {
		return errors.New("socket is not connected")
	}
	go c.listener()
	return nil
}

func (c *PubSubClient) listener() {
	for c.isConnected {
		data, err := c.Socket.RecvBytes(0)
		if err != nil {
			panic(err)
		}
		message := pubSubFrame{}
		err = json.Unmarshal(data, &message)
		if err != nil {
			panic(err)
		}
		c.topicHandlers[message.Topic](message.Data)
	}
}
