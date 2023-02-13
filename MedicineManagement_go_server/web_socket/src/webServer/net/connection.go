package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/utils"

	"github.com/gorilla/websocket"
)

// Connection 当前链模块
type Connection struct {
	//当前conn时属于哪个server的
	WebServer iface.IServer
	// 当前链接的socket
	Conn *websocket.Conn
	//ID
	ConnID uint32
	//当前的链接状态
	ISClosed bool

	//告知当前链接已经停止/退出的channel,由reader告知writer退出
	ExitChan chan bool

	// 无缓冲的管道，用户读写goroutine之间的消息通信
	msgChan chan []byte

	// 当前请求的request
	Req *http.Request

	// 消息的管理MsgID和对应的处理API关系
	MsgHandle iface.IMsgHandle

	// 链接属性集合
	Property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接模块的方法
func NewConnection(server iface.IServer, conn *websocket.Conn, connID uint32, msgHandle iface.IMsgHandle, req *http.Request) *Connection {
	c := &Connection{
		WebServer: server,
		Conn:      conn,
		ConnID:    connID,
		ISClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
		MsgHandle: msgHandle,
		Req:       req,
		Property:  make(map[string]interface{}),
	}

	// 将conn加入到connManager中
	c.WebServer.GetConnMgr().Add(c)

	return c
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")

	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		_, mb, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("ReadMsg Error", err)
			return
		}
		msg := &Message{}
		dataBuff := bytes.NewReader(mb)
		// 先读ID
		err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
		if err != nil {
			return
		}
		// 在读data
		msg.Data, err = io.ReadAll(dataBuff)
		if err != nil {
			return
		}

		// 得到当前连接的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制
			// 根据绑定好的msgID找到对应的处理api业务执行
			c.MsgHandle.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandle.DoMsgHandler(req)
		}
	}
}

// StartWriter 写消息goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, "Writer is exit, remote addr is ", c.RemoteAddr())
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				fmt.Println("Send Data Error ", err)
				continue
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)
	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接的写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的创建链接后需要调用的处理业务
	c.WebServer.CallOnConnStart(c)
}
func (c *Connection) Stop() {
	fmt.Println("Conn Stop.. ConnID = ", c.ConnID)
	// 如果当前链接已经关闭
	if c.ISClosed {
		return
	}

	// 调用开发者注册的销毁链接之前需要执行的hook函数
	c.WebServer.CallOnConnStop(c)

	c.ISClosed = true
	//关闭socket链接
	err := c.Conn.Close()
	if err != nil {
		return
	}

	// 将当前链接从connMgr中摘除掉
	c.WebServer.GetConnMgr().Remove(c)

	// 告知Writer关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)

}
func (c *Connection) GetConnection() *websocket.Conn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() string {
	return c.Req.RemoteAddr
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.ISClosed {
		return errors.New("connection is closed")
	}
	// 将数据发送给客户端
	c.msgChan <- data
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.Property[key] = value
	fmt.Println("AddProperty ", key, " is ", value)
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if v, ok := c.Property[key]; ok {
		return v, nil
	}
	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.Property, key)
}

// GetReq 获取当前连接的请求
func (c *Connection) GetReq() *http.Request {
	return c.Req
}
