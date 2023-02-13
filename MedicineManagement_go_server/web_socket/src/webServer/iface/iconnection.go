package iface

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// IConnection 定义链接模块的抽象类
type IConnection interface {
	// Start 启动链接 让当前的链接准备开始工作
	Start()
	// Stop 停止链接 结束当前链接的工作
	Stop()
	// GetConnection 获取当前链接绑定的socket conn
	GetConnection() *websocket.Conn
	// GetReq 获取当前连接的请求
	GetReq() *http.Request
	// GetConnID 获取当前链接模块的链接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的TCP状态，IP，Port
	RemoteAddr() string
	// SendMsg 发送数据 将数据发送给远程的客户端
	SendMsg(id uint32, data []byte) error
	// SetProperty 设置链接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取链接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除链接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*websocket.Conn, []byte) error
