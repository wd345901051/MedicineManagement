package core

import (
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/webServer/iface"
)

type RoomMst struct {
	Info *define.JwtClaims
	Conn iface.IConnection
}

func NewRoomMst(claim *define.JwtClaims, conn iface.IConnection) *RoomMst {
	return &RoomMst{
		Info: claim,
		Conn: conn,
	}
}
