package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

// User 用户结构体
type User struct {
	conn *websocket.Conn
	msg  chan []byte
}
type Msg struct {
	ID  uint32 `json:"id"`
	Msg []byte `json:"msg"`
}

func read(user *User) {
	// 从用户chan中读取数据
	for {
		_, msg, err := user.conn.ReadMessage()
		if err != nil {
			fmt.Println("用户退出：", user.conn.RemoteAddr().String())
			hb.unregister <- user
			break
		}
		fmt.Println("读到数据啦！！！", msg)
		dataBuff := bytes.NewReader(msg)

		m := &Msg{}
		fmt.Println("ReadBuff", dataBuff)
		// 先读ID
		if err := binary.Read(dataBuff, binary.LittleEndian, &m.ID); err != nil {
			return
		}
		fmt.Println("我还省点", dataBuff)
		// 再读msg
		all, err2 := io.ReadAll(dataBuff)
		if err2 != nil {
			return
		}
		fmt.Println("Msg Is ", string(all))
		// 将读到的消息进行广播
		fmt.Println(m, "Read Msg Success!!")

		hb.broadcast <- msg
	}
}

func wirte(user *User) {
	for data := range user.msg {
		err := user.conn.WriteMessage(1, data)
		if err != nil {
			fmt.Println("User Write Msg Error:", err)
			break
		}
	}
}

// Hub 消息中心
type Hub struct {
	// 用户列表
	userList map[*User]bool
	// 注册chan，用户注册是添加到chan中
	register chan *User
	// 注销chan，用户退出是添加到chan中，再从map中删除
	unregister chan *User
	//广播消息，将消息广播给所有连接
	broadcast chan []byte
}

// 消息中心的消息处理
func (h *Hub) run() {
	for {
		select {
		//注册
		case user := <-h.register:
			h.userList[user] = true
		case user := <-h.unregister:
			if _, ok := h.userList[user]; ok {
				delete(h.userList, user)
				close(user.msg)
				user.conn.Close()
			}
		case data := <-h.broadcast:
			//从广播消息chan中取出消息，并便利用户发送
			for u := range h.userList {
				u.msg <- data
			}
		}
	}
}

var up = &websocket.Upgrader{
	// 定义读写缓冲区
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		// 如果不是get请求，返回错误
		if r.Method != "GET" {
			fmt.Println("请求方式错误!")
			return false
		}
		//如果路径中不包括socket，返回错误
		if r.URL.Path != "/socket" {
			fmt.Println("请求路径错误")
			return false
		}
		// 还可以自定义规则
		return true
	},
}

var hb = &Hub{
	userList:   make(map[*User]bool),
	register:   make(chan *User),
	unregister: make(chan *User),
	broadcast:  make(chan []byte),
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	// 通过审计后的升级器得到连接，将HTTP请求升级为WebSocket连接
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("获取连接失败")
		return
	}
	// 连接之后的手续
	user := &User{
		conn: conn,
		msg:  make(chan []byte),
	}
	// 将用户注册进消息队列
	hb.register <- user
	defer func() {
		hb.unregister <- user
	}()
	// TODO读到连接后处理用户业务
	go read(user)
	wirte(user)
}

func main() {
	// 后台启动处理器
	go hb.run()
	http.HandleFunc("/socket", wsHandle)
	fmt.Println("Server Is Start...")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}
}
