package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Msg struct {
	ID   uint32 `json:"id"`
	Data string `json:"msg"`
}

func main() {
	r := http.Header{}
	var au string
	fmt.Println("Place Input Au")
	_, err := fmt.Scanf("%s", &au)
	if err != nil {
		return
	}
	r.Add("Authorization", au)
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8999/socket", r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("客户端启动成功")
	go Read(conn)
	for true {
		m := &Msg{}
		fmt.Print("输入msgID,msgData:\n")
		_, err := fmt.Scan(&m.ID)
		_, err = fmt.Scan(&m.Data)
		if err != nil {
			fmt.Println(err)
			return
		}

		dataBuff := bytes.NewBuffer([]byte{})
		// 先写入ID
		if err := binary.Write(dataBuff, binary.LittleEndian, m.ID); err != nil {
			fmt.Println("Write Msg ID Error")
			return
		}
		// 在写入data
		fmt.Println(m.Data)
		if err := binary.Write(dataBuff, binary.LittleEndian, []byte(m.Data)); err != nil {
			fmt.Println("Write Msg Data Error")
			return
		}
		err = conn.WriteMessage(websocket.BinaryMessage, dataBuff.Bytes())
		if err != nil {
			fmt.Println("Send Msg Error")
			return
		}
		fmt.Println("Send Msg Succ!!", "msgID:", m.ID, "msgData", m.Data)
	}
}

func Read(conn *websocket.Conn) {
	for {
		_, mb, err := conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println("服务器发来消息：", string(mb))
	}
}
