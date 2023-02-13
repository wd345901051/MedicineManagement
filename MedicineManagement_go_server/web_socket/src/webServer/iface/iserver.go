package iface

import "net/http"

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务
	Start(http.ResponseWriter, *http.Request)
	// Stop 停止服务
	Stop()
	// Server 运行服务
	Server()
	// AddRouter 路由功能。给当前的服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前的ConnMgr
	GetConnMgr() IConnManager
	// SetOnConnStart 注册OnConnStart方法
	SetOnConnStart(func(connection IConnection))
	// SetOnConnStop 注册OnConnStop方法
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStart 调用OnConnStart方法
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用OnConnStop方法
	CallOnConnStop(connection IConnection)
}
