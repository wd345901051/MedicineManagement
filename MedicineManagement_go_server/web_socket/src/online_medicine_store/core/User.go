package core

import (
	"strconv"
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/online_medicine_store/models"
	"web_socket/src/webServer/iface"
)

type User struct {
	Info   *define.JwtClaims
	Conn   iface.IConnection
	Master *RoomMst
	RoomID uint32
}

func NewUser(claims *define.JwtClaims, conn iface.IConnection, mst *RoomMst) *User {
	models.RDB.Incr(define.REDIS_ROOM_TALK_USER_NUMBER)
	return &User{
		Info:   claims,
		Conn:   conn,
		Master: mst,
	}
}

func (u *User) UserLoginSendHistoryMsg(msgs []string) {
	for _, v := range msgs {
		_ = u.Conn.SendMsg(1, []byte(v))
	}
	models.RDB.Del(define.REDIS_USER_TALK_QUEUE + strconv.Itoa(int(u.Info.ID)))
}
