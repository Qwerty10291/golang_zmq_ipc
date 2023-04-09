package objects

import "github.com/pebbe/zmq4"

type SocketType int16

const (
	REP_SERVER    SocketType = 4
	PUB_SERVER    SocketType = 1
	PUSH_SERVER   SocketType = 8
	ROUTER_SERVER SocketType = 6
	DEALER_SERVER SocketType = 5
)

type App struct {
	Name        string   `json:"name"`
	Initialized bool     `json:"initialized"`
	MainServer  Server   `json:"main_server"`
	Servers     []Server `json:"servers"`
	MetricsPort *int     `json:"metrics_port"`
}

type Server struct {
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	SocketType int    `json:"socket"`
}

type Client interface {
	Connect() error
	Close() error
	GetSocket() *zmq4.Socket
}

type ServerInterface interface {
	Bind() error
	Close() error
	GetSocket() *zmq4.Socket
}
