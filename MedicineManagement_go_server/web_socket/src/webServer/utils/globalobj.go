package utils

import (
	"encoding/json"
	"io/ioutil"
	"web_socket/src/webServer/iface"
)

// GlobalObj 存储一切有关框架的全局参数，供其他模块使用，一些参数时可以通过json由用户进行配置的
type GlobalObj struct {
	Server           iface.IServer // 当前服务的全局对象
	Host             string        // 当前服务器主机监听的IP
	Port             int           // 当前服务器的端口
	Name             string        // 当前服务器的名称
	Version          string        // 当前服务的版本
	MaxConn          int           // 允许的最大连接数
	MaxPackageSize   int           // 当前框架允许的最大数据包值
	WorkerPoolSize   uint32        // WorkerGoroutine数量
	MaxWorkerTaskLen uint32        // 框架允许的最大Task数量
}

// GlobalObject 定义一个全局的对外Global对象
var GlobalObject *GlobalObj

// Reload 从配置文件中加载配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject对象
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		Port:             8999,
		Name:             "WebServerApp",
		Version:          "V0.9",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	// 应该尝试从conf/server.json取加载一些用户自定义的参数
	GlobalObject.Reload()
}
