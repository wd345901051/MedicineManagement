package apis

import (
	"encoding/json"
	"web_socket/src/online_medicine_store/core"
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/net"
)

type TalkToRoomMember struct {
	net.BaseRouter
}

/*   给房间内某人发送消息,必须时房主
{
    "id":2,
    "data":"{
		"id",1,
		"msg":"哈哈"
	}"
}
*/

func (rt *TalkToRoomMember) Handle(req iface.IRequest) {
	iRoom, err := req.GetConnection().GetProperty("Room")
	if err != nil {
		return
	}
	room, ok := iRoom.(*core.Room)
	if !ok {
		return
	}
	if room.Master.Info.Root != 1 {
		return
	}
	info := &define.Action2{}
	err = json.Unmarshal(req.GetMsg().GetData(), info)
	if err != nil {
		return
	}
	room.SendMsgToMember(info.Id, []byte(info.Msg))
}
