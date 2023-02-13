package apis

import (
	"web_socket/src/online_medicine_store/core"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/net"
)

type TalkToRoomMaster struct {
	net.BaseRouter
}

/*   给房主发送消息
{
    "id":1,
    "data":"哈哈"
}
*/

func (rt *TalkToRoomMaster) Handle(req iface.IRequest) {
	iUs, err := req.GetConnection().GetProperty("User")
	if err != nil {
		return
	}
	us, ok := iUs.(*core.User)
	if !ok {
		return
	}
	_ = us.Master.Conn.SendMsg(1, req.GetMsg().GetData())
}
