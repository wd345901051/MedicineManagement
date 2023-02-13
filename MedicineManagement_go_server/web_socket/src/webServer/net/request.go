package net

import "web_socket/src/webServer/iface"

type Request struct {
	// 已经和客户端建立好连接的conn
	conn iface.IConnection
	// 客户端请求的数据
	msg iface.IMessage
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// GetMsg 得到当前数据
func (r *Request) GetMsg() iface.IMessage {
	return r.msg
}
