package iface

type IConnManager interface {
	// Add 添加链接
	Add(conn IConnection)
	// Remove 删除链接
	Remove(conn IConnection)
	// Get 更具connID获取链接
	Get(connID uint32) (IConnection, error)
	// Len 得到道歉链接总数
	Len() int
	// ClearConn 清除并中止所有链接
	ClearConn()
}
