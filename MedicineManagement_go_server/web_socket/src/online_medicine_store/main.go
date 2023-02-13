package main

import (
	"fmt"
	"strconv"
	"web_socket/src/online_medicine_store/apis"
	"web_socket/src/online_medicine_store/core"
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/online_medicine_store/models"
	"web_socket/src/online_medicine_store/utils"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/net"
)

func ConnStart(conn iface.IConnection) {
	token := conn.GetReq().Header.Get("Authorization")
	claim, err := utils.AnalysisToken(token)
	if err != nil {
		err = conn.SendMsg(1, []byte("Token Error"))
		conn.Stop()
		return
	}

	if claim.Root == 1 { //高级用户
		mst := core.NewRoomMst(claim, conn)
		room := core.NewRoom(mst)
		core.RoomMgrObj.AddRoom(room)
		_ = conn.SendMsg(1, []byte(fmt.Sprint(claim.Name, "创建了", room.ID, "号房间")))
		conn.SetProperty("Room", room)
	}
	if claim.Root == 0 { //普通用户
		rn, _ := models.RDB.Get(define.REDIS_ROOM_COUNT).Int()
		if rn < 1 {
			_ = conn.SendMsg(1, []byte("当前不存在房间"))
			conn.Stop()
			return
		}
		room := core.RoomMgrObj.GetOneRoom()
		us := core.NewUser(claim, conn, room.Master)
		room.AddMember(claim.ID, us)
		conn.SetProperty("User", us)
		_ = conn.SendMsg(1, []byte(fmt.Sprint(claim.Name, "加入了", room.ID, "号房间")))
		result, _ := models.RDB.LRange(define.REDIS_USER_TALK_QUEUE+strconv.Itoa(int(us.Info.ID)), 0, -1).Result()
		go us.UserLoginSendHistoryMsg(result)
	}
	conn.SetProperty("Claim", claim)
}

func ConnStop(conn iface.IConnection) {
	// 回收用户资源
	iClaim, err := conn.GetProperty("Claim")
	if err != nil {
		return
	}
	claim, ok := iClaim.(*define.JwtClaims)
	if !ok {
		return
	}
	if claim.Root == 0 {
		iUser, err := conn.GetProperty("User")
		if err != nil {
			return
		}
		us, ok := iUser.(*core.User)
		if !ok {
			return
		}
		room := core.RoomMgrObj.GetRoomById(us.RoomID)
		room.RemoveMember(us.Info.ID)
		_ = conn.SendMsg(1, []byte(fmt.Sprint(claim.Name, "离开了", room.ID, "号房间")))
		models.RDB.Decr(define.REDIS_ROOM_TALK_USER_NUMBER)
	}
	if claim.Root == 1 {
		iRoom, err := conn.GetProperty("Room")
		if err != nil {
			return
		}
		room, ok := iRoom.(*core.Room)
		if !ok {
			return
		}
		_ = core.RoomMgrObj.RemoveRoom(room.ID)
	}
}

func main() {
	server := net.NewServer("WebServer")
	server.SetOnConnStart(ConnStart)
	server.SetOnConnStop(ConnStop)
	// 注册路由
	server.AddRouter(1, &apis.TalkToRoomMaster{})
	server.AddRouter(2, &apis.TalkToRoomMember{})
	server.AddRouter(3, &apis.TalkToRoomAll{})
	server.Server()
}
