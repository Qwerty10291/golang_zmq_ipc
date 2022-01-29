package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Qwerty10291/golang_zmq_ipc/client"
	"github.com/Qwerty10291/golang_zmq_ipc/objects"
	"github.com/Qwerty10291/golang_zmq_ipc/server"
	"github.com/pebbe/zmq4"
)

type App struct {
	Name      string
	Servers   map[string]objects.Server
	OtherApps map[string]objects.App

	controllerClient *client.ReqRepClient
	controllerServer *server.ReqRepServer

	config     AppConfig
	zmqContext *zmq4.Context
	logger     *log.Logger
}

func NewApp(appName string, config AppConfig, zmqContext *zmq4.Context, logger *log.Logger) *App {
	return &App{
		Name:             appName,
		Servers:          map[string]objects.Server{},
		OtherApps:        map[string]objects.App{},
		controllerClient: &client.ReqRepClient{},
		controllerServer: &server.ReqRepServer{},
		config:           config,
		zmqContext:       zmqContext,
		logger:           logger,
	}
}

type AppConfig struct {
	ControllerHost string
	ControllerPort string
	// milliseconds
	ControllerResponseTimeout time.Duration
}

func (a *App) Init() error {
	var err error
	a.controllerClient, err = client.NewReqRepClient(a.config.ControllerHost, a.config.ControllerPort, client.TCP, a.zmqContext, a.config.ControllerResponseTimeout)
	if err != nil {
		return err
	}
	err = a.controllerClient.Connect()
	if err != nil{
		log.Println("error when connecting controller host")
		return err
	}
	a.logger.Printf("controller client created at %s:%s", a.controllerClient.Host, a.controllerClient.Port)
	err = a.WaitForControllerInit()
	if err != nil {
		log.Println("error when waiting for controller init")
		return err
	}
	err = a.register()
	if err != nil {
		return err
	}
	a.logger.Println("initialized")
	return nil
}

func (a *App) WaitForControllerInit() error {
	for {
		resp, err := a.controllerClient.RequestRaw("ping", struct{}{})
		if err != nil {
			return err
		}
		var status bool
		err = json.Unmarshal(resp, &status)
		if err != nil {
			return err
		}
		if status {
			break
		}
	}
	return nil
}

func (a *App) NewServer(serverName string, socketType objects.SocketType, protocol server.Protocol, context *zmq4.Context) (objects.ServerInterface, error) {
	if server, ok := a.Servers[serverName]; ok {
		return a.createServer(server, context)
	} else {
		data, err := a.controllerClient.RequestRaw("register_server", controllerRegisterServerRequest{
			AppName:    a.Name,
			ServerName: serverName,
			SocketType: socketType,
			Protocol:   protocol,
		})
		if err != nil {
			return nil, err
		}
		var serverData objects.Server
		err = json.Unmarshal(data, &serverData)
		if err != nil {
			return nil, err
		}
		server, err := a.createServer(serverData, context)
		if err != nil {
			return nil, err
		}
		a.Servers[serverName] = serverData
		return server, nil
	}
}

func (a *App) NewClient(appName string, serverName string, context *zmq4.Context) (objects.Client, error) {
	if app, ok := a.OtherApps[appName]; ok {
		for _, server := range app.Servers {
			if server.Name == serverName {
				return a.createClient(server, context)
			}
		}
		return nil, fmt.Errorf("sever with name %s not found in %s", serverName, appName)
	} else {
		err := a.updateOtherApps()
		if err != nil {
			return nil, err
		}
		if app, ok := a.OtherApps[appName]; ok {
			for _, server := range app.Servers {
				if server.Name == serverName {
					return a.createClient(server, context)
				}
			}
			return nil, fmt.Errorf("sever with name %s not found in %s", serverName, appName)
		} else {
			return nil, fmt.Errorf("app with name %s not found", appName)
		}
	}
}

func (a *App) LoadOtherAppsData(local bool) (map[string]objects.App, error) {
	data, err := a.controllerClient.RequestRaw("get_apps", struct {
		Local bool `json:"local"`
	}{Local: local})
	if err != nil {
		return nil, err
	}
	var apps map[string]objects.App
	err = json.Unmarshal(data, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (a *App) WaitForApp(appName string) error {
	for {
		err := a.updateOtherApps()
		if err != nil {
			return err
		}
		if _, ok := a.OtherApps[appName]; ok {
			return nil
		}
	}
}

func (a *App) register() error {
	response, err := a.controllerClient.RequestRaw("register_app", struct {
		Name string `json:"name"`
	}{a.Name})
	if err != nil {
		return err
	}
	appInfo := serverRegisterAppResponse{}
	err = json.Unmarshal(response, &appInfo)
	if err != nil || appInfo.Port == 0 {
		alreadyExistsData := controllerAppAlreadyExistsError{}
		err = json.Unmarshal(response, &alreadyExistsData)
		if err != nil {
			a.logger.Fatalf("unknown controller response type:%s", string(response))
		}
		a.logger.Println("app already exists, restore")
		return a.restoreFromControllerResponse(alreadyExistsData.App)
	} else {
		a.logger.Println("app does nor existed")
		a.controllerServer, err = a.initControllerServer(appInfo)
		if err != nil {
			return err
		}
		a.logger.Printf("controller server created at %s:%s\n", a.controllerServer.Host, a.controllerServer.Port)
	}
	return nil
}

func (a *App) initControllerServer(resp serverRegisterAppResponse) (*server.ReqRepServer, error) {
	port := strconv.Itoa(resp.Port)
	var err error
	server, err := server.NewReqRepServer("*", port, server.TCP, a.zmqContext)
	if err != nil{
		return nil, err
	}
	err = server.Bind()
	if err != nil{
		return nil, err
	}
	server.Start()
	return server, nil
}

func (a *App) updateOtherApps() error {
	apps, err := a.LoadOtherAppsData(false)
	if err != nil {
		a.logger.Println("error when updating other apps data")
		return err
	}
	a.OtherApps = apps
	return nil
}

func (a *App) restoreFromControllerResponse(app objects.App) error {
	var err error
	a.controllerServer, err = server.NewReqRepServer("*", strconv.Itoa(app.MainServer.Port), server.Protocol(app.MainServer.Protocol), a.zmqContext)
	if err != nil {
		return err
	}
	for _, serverData := range app.Servers {
		a.Servers[serverData.Name] = serverData
	}
	return nil
}

func (a *App) createServer(serverData objects.Server, context *zmq4.Context) (objects.ServerInterface, error) {
	switch serverData.SocketType {
	case int(objects.PUB_SERVER):
		a.logger.Printf("new pub_sub server: ip:%s port:%d protocol:%s", serverData.Ip, serverData.Port, serverData.Protocol)
		return server.NewPubSubServer("*", strconv.Itoa(serverData.Port), server.Protocol(serverData.Protocol), a.zmqContext)
	case int(objects.REP_SERVER):
		a.logger.Printf("new pub_sub server: ip:%s port:%d protocol:%s", serverData.Ip, serverData.Port, serverData.Protocol)
		return server.NewPubSubServer("*", strconv.Itoa(serverData.Port), server.Protocol(serverData.Protocol), a.zmqContext)
	default:
		return nil, fmt.Errorf("socket with id %d not found", serverData.SocketType)
	}
}

func (a *App) createClient(serverData objects.Server, context *zmq4.Context) (objects.Client, error) {
	switch serverData.SocketType {
	case int(objects.PUB_SERVER):
		a.logger.Printf("new pub_sub client: ip:%s port:%d protocol:%s", serverData.Ip, serverData.Port, serverData.Protocol)
		return client.NewPubSubClient(serverData.Ip, strconv.Itoa(serverData.Port), client.Protocol(serverData.Protocol), context)
	case int(objects.REP_SERVER):
		a.logger.Printf("new req_rep client: ip:%s port:%d protocol:%s", serverData.Ip, serverData.Port, serverData.Protocol)
		return client.NewReqRepClient(serverData.Ip, strconv.Itoa(serverData.Port), client.Protocol(serverData.Protocol), context, 0)
	default:
		return nil, fmt.Errorf("socket with id %d not found", serverData.SocketType)
	}
}
