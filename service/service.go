package service

import (
	"strconv"
	"sync"
	"time"
)

type ProxyServer struct {
	Host       string    `json:"host"`
	Port       int       `json:"port"`
	LastAction time.Time `json:"lastAction"`
	IsAction   bool      `json:"isAction"`

	Mu sync.Mutex
}

func NewProxyServerByModel(model ProxyServerModel) *ProxyServer {
	return &ProxyServer{
		Host: model.Host,
		Port: model.Port,
	}
}

type ProxyServerModel struct {
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	RequestTime time.Time `json:"requestTime"`
}


func (p *ProxyServer) GetKey() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}