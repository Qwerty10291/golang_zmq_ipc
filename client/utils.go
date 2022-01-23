package client

import zmq "github.com/pebbe/zmq4"

type Protocol string 

const (
	TCP Protocol = "tcp" 
)


func SetSocketFlags(socket *zmq.Socket) {
	socket.SetTcpKeepalive(1)
	socket.SetTcpKeepaliveIdle(30)
	socket.SetTcpKeepaliveCnt(5)
	socket.SetTcpKeepaliveIntvl(5)
	socket.SetLinger(0)
}