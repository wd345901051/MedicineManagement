package test

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"
)

type Msg struct {
	Id  uint32
	Msg []byte
}

func TestTranslate(t *testing.T) {
	// Create the claims
	m := Msg{
		Id:  2,
		Msg: []byte(`{"id":2,"msg":"你下线啦？2"}`),
	}
	buffer := bytes.NewBuffer([]byte{})
	err := binary.Write(buffer, binary.LittleEndian, m.Id)
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, m.Msg)
	if err != nil {
		return
	}
	encode := base64.StdEncoding.EncodeToString(buffer.Bytes())
	fmt.Println(encode)
}
