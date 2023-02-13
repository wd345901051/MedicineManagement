package apis

import (
	"web_socket/src/online_medicine_store/core"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/net"
)

type TalkToRoomAll struct {
	net.BaseRouter
}

/*   给房间内所有人发送消息，必须是房主
{
    "id":3,
    "data":"哈哈"
}
*/

func (rt *TalkToRoomAll) Handle(req iface.IRequest) {
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
	err = room.SendMsgToAll(req.GetMsg().GetData())
	if err != nil {
		return
	}
}
