package app

import "github.com/Qwerty10291/golang_zmq_ipc/objects"

type controllerAppAlreadyExistsError struct{
	Error int `json:"error"`
	App objects.App `json:"app"`
}