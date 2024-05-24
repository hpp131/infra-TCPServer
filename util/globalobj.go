package util

import (
	"encoding/json"
	"os"
	"tcpserver/ziface"
)

func init() {
	Globalobject = &GlobalObj{
		Version:        "v1",
		Name:           "TcpServerv1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		Host:           "0.0.0.0",
		TCPPort:        8999,
		MaxWorkTaskLen: 1024,
		WorkPoolSize:   10,
	}
	Globalobject.Load()
}

type GlobalObj struct {
	TCPServer ziface.IServer
	TCPPort   uint32
	// 当前server的版本号，并非TCPVersion
	Version        string
	Host           string
	Name           string
	MaxConn        int
	MaxPackageSize uint32
	WorkPoolSize   uint32
	// 每个worker所消费的TaskQueue的最大长度
	MaxWorkTaskLen uint32
}

var Globalobject *GlobalObj

// 从配置文件加载一部分参数
func (g *GlobalObj) Load() {
	data, err := os.ReadFile("conf/server.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &Globalobject)
	if err != nil {
		panic(err)
	}
}
