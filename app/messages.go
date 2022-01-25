package app

import (
	"github.com/Qwerty10291/golang_zmq_ipc/objects"
	"github.com/Qwerty10291/golang_zmq_ipc/server"
)

type serverRegisterAppResponse struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type controllerRegisterServerRequest struct{
	AppName string `json:"app_name"`
	ServerName string `json:"server_name"`
	SocketType objects.SocketType `json:"socket_type"`
	Protocol server.Protocol `json:"protocol"`
}