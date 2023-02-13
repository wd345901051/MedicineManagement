package test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"strconv"
	"testing"
	"web_socket/src/online_medicine_store/define"
)

const JwtKey = "haha"

type JwtClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func TestGenerateToken(t *testing.T) {
	// Create the claims
	claims := define.JwtClaims{
		ID:               4,
		Name:             "赵六",
		Root:             0,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return
	}
	fmt.Println(ss)
}

func TestAnalysisToken(t *testing.T) {
	// Create the claims
	tokens := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjpbeyJuYW1lIjoi5byg5LiJIn1dLCJpYXQiOjE2Njc3MDg3MzEsImV4cCI6MTY2ODk1OTk5OSwiYXVkIjoiIiwiaXNzIjoiIiwic3ViIjoiIn0.iy3-KA1P2UPFrjp6Rjs1qjHseQi2aa0a4egk7Np7dPo"
	claim := JwtClaims{}
	token, err := jwt.ParseWithClaims(tokens, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		t.Fatal("Analysis Error", err)
		return
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		fmt.Println(claims.Name)
	} else {
		t.Fatal("Error")
		return
	}
}

type A struct {
	Id  uint32 `json:"id"`
	Msg string `json:"msg"`
}

type B struct {
	Id  uint32 `json:"id"`
	Msg string `json:"msg"`
}

func Test123(t *testing.T) {
	b := B{
		Id:  18,
		Msg: "haha",
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return
	}
	br := bytes.NewReader(bb)

	header1 := make([]byte, 6)
	err = binary.Read(br, binary.LittleEndian, header1) // 读取 {"id":
	header2 := make([]byte, 2)                          // 读取id
	_, err = io.ReadFull(br, header2)
	header3 := make([]byte, 8) // 读取,msg:
	_, err = io.ReadFull(br, header3)
	header4, err := io.ReadAll(br)
	//s:=[]byte(`{"id":1,msg:"123"}`)

	fmt.Println(string(header2), string(header4[:len(header4)-2]))
}

func Test456(t *testing.T) {
	fmt.Println(strconv.Atoi("001"))
}
