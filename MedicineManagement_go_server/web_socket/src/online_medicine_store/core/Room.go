package core

import (
	"fmt"
	"strconv"
	"sync"
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/online_medicine_store/models"
)

type Room struct {
	ID       uint32
	Master   *RoomMst
	Members  map[uint32]*User // 当前房间所有人的name链接conn（不包含主人）
	RoomLock sync.RWMutex
}

// NewRoom 创建一个房间的方法
func NewRoom(master *RoomMst) *Room {
	//TODO ID = sql.masterID
	roomId, err := models.RDB.Get(define.REDIS_ROOM_COUNT).Int()
	if err != nil {
		return nil
	}
	i := 0
	for i < roomId {
		if _, ok := RoomMgrObj.Rooms[uint32(i)]; !ok {
			roomId = i
			break
		}
		i++
	}
	if i == roomId {
		models.RDB.Incr(define.REDIS_ROOM_COUNT)
	}
	fmt.Println("创建了", roomId, "号房间----------------------")
	return &Room{
		ID:      uint32(roomId),
		Master:  master,
		Members: make(map[uint32]*User),
	}
}

// AddMember 添加一个成员的方法
func (r *Room) AddMember(id uint32, member *User) {
	r.RoomLock.Lock()
	defer r.RoomLock.Unlock()
	if _, ok := r.Members[id]; ok {
		fmt.Println("This Member Is Exist")
		return
	}
	r.Members[id] = member
	member.RoomID = r.ID
	err := r.SendMsgToMaster([]byte(r.Members[id].Info.Name + "进入了本房间"))
	if err != nil {
		return
	}
}

// RemoveMember 移除一个成员的方法
func (r *Room) RemoveMember(id uint32) {
	r.RoomLock.Lock()
	defer r.RoomLock.Unlock()
	if _, ok := r.Members[id]; !ok {
		fmt.Println("This Member Is Not Exist")
		return
	}
	err := r.SendMsgToMaster([]byte(r.Members[id].Info.Name + "离开了本房间"))
	if err != nil {
		return
	}
	delete(r.Members, id)
}

func (r *Room) SendMsgToMaster(data []byte) error {
	err := r.Master.Conn.SendMsg(1, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Room) SendMsgToMember(id uint32, data []byte) {
	r.RoomLock.RLock()
	defer r.RoomLock.RUnlock()
	if member, ok := r.Members[id]; ok {
		_ = member.Conn.SendMsg(1, data)
	} else {
		lenght, _ := models.RDB.LLen(define.REDIS_USER_TALK_QUEUE + strconv.Itoa(int(id))).Result()
		if lenght > 20 {
			// TODO 限制存储的最大聊天记录数
			return
		}
		models.RDB.RPush(define.REDIS_USER_TALK_QUEUE+strconv.Itoa(int(id)), data)
	}
}

func (r *Room) MasterLeaved() error {
	for _, usConn := range r.Members {
		_ = usConn.Conn.SendMsg(1, []byte("房主已下线"))
		usConn.Conn.Stop()
	}
	return nil
}

func (r *Room) SendMsgToAll(data []byte) error {
	for _, usConn := range r.Members {
		_ = usConn.Conn.SendMsg(1, data)
	}
	return nil
}
