package iface

// IRequest 实际上是把客户端请求的连接信息和请求的数据包装到一个request中
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetMsg 得到当前msg
	GetMsg() IMessage
}
