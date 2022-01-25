package objects

type SocketType int16

const (
	REP_SERVER  SocketType = 4
	PUB_SERVER  SocketType = 1
	PUSH_SERVER SocketType = 8
)

type App struct {
	Name        string   `json:"name"`
	Initialized bool     `json:"initialized"`
	MainServer  Server   `json:"main_server"`
	Servers     []Server `json:"servers"`
}

type Server struct {
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Port       int `json:"port"`
	Protocol   string `json:"protocol"`
	SocketType int    `json:"socket"`
}

type Client interface {
	Connect() error
	Close() error
}

type ServerInterface interface {
	Bind() error
	Close() error
}