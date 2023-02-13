package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/utils"
)

// Server 的接口实现，定义一个Server服务模块
type Server struct {
	// 服务器的名称
	Name string

	// 服务器所监听的IP
	IP string
	// 服务器所监听的端口
	Port int

	// websocket配置服务
	Wup *websocket.Upgrader

	// 当前server消息管理模块，用来绑定msgID和对应的处理关系
	MsgHandler iface.IMsgHandle

	// 该server的链接管理器
	ConnMgr iface.IConnManager

	// 创建conn后的Hook函数
	OnConnStart func(conn iface.IConnection)
	// 销毁conn前的Hook函数
	OnConnStop func(conn iface.IConnection)
}

var cid uint32

// Start 启动服务
func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	// 从http请求中获取connection
	conn, err := s.Wup.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	}
	// 设置最大链接个数的判断，如果超过最大链接的数量，那么则关闭此新的链接
	if s.ConnMgr.Len() > utils.GlobalObject.MaxConn {
		fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	// 将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
	dealConn := NewConnection(s, conn, cid, s.MsgHandler, r)
	cid++

	// 启动当前的链接业务处理
	go dealConn.Start()
}

// Stop 停止服务
func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些已经开辟的连接信息停止或者回收
	fmt.Println("[STOP] Server ", s.Name)
	s.ConnMgr.ClearConn()
}

// Server 运行服务
func (s *Server) Server() {
	// 启动server的服务功能
	// 开启消息队列及worker工作吃
	s.MsgHandler.StartWorkerPool()

	http.HandleFunc("/", s.Start)
	fmt.Printf("[WebServer] Server Name : %s, listenner at IP : %s, Port : %d is starting \n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.Port)

	fmt.Printf("[WebServer] Version : %s ,MaxConn : %d,MaxPackageSize : %d \n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	//TODO 做一些启动服务之后的业务
	err := http.ListenAndServe(fmt.Sprint(s.IP, ":", s.Port), nil)
	if err != nil {
		fmt.Println("Start HTTP Server Error:", err)
		return
	}
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!!")
}

// NewServer 初始化server模块的方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IP:   utils.GlobalObject.Host,
		Port: utils.GlobalObject.Port,
		Wup: &websocket.Upgrader{
			// 定义读写缓冲区
			WriteBufferSize: utils.GlobalObject.MaxPackageSize,
			ReadBufferSize:  utils.GlobalObject.MaxPackageSize,
			CheckOrigin: func(r *http.Request) bool {
				// 如果不是get请求，返回错误
				if r.Method != "GET" {
					fmt.Println("请求方式错误!")
					return false
				}
				////如果路径中不包括socket，返回错误
				//if r.URL.Path != "/socket" {
				//	fmt.Println("请求路径错误")
				//	return false
				//}
				// 还可以自定义规则
				return true
			},
		},
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 注册OnConnStart方法
func (s *Server) SetOnConnStart(hookFunc func(connection iface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop方法
func (s *Server) SetOnConnStop(hookFunc func(connection iface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart方法
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop方法
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("------> Call OnConnStart()...")
		s.OnConnStop(conn)
	}
}
