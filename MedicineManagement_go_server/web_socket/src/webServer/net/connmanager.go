package net

import (
	"errors"
	"fmt"
	"sync"
	"web_socket/src/webServer/iface"
)

type ConnManager struct {
	Connections map[uint32]iface.IConnection //管理的链接信息集合
	connLock    sync.RWMutex                 // 保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]iface.IConnection),
	}
}

// Add 添加链接
func (cm *ConnManager) Add(conn iface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.Connections[conn.GetConnID()] = conn
	fmt.Println("ConnID = ", conn.GetConnID(), " add to ConnManager Successfully: conn num = ", cm.Len())
}

// Remove 删除链接
func (cm *ConnManager) Remove(conn iface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.Connections, conn.GetConnID())
	fmt.Println("ConnID = ", conn.GetConnID(), " remove to ConnManager Successfully: conn num = ", cm.Len())
}

// Get 更具connID获取链接
func (cm *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.Connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not FOUND")
}

// Len 得到道歉链接总数
func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

// ClearConn 清除并中止所有链接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除conn并停止conn的工作
	for connID, conn := range cm.Connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.Connections, connID)
	}
	fmt.Println("Clear All Connections Succ! conn num = ", cm.Len())
}
