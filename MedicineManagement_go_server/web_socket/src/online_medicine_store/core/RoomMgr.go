package core

import (
	"fmt"
	"sync"
	"web_socket/src/online_medicine_store/define"
	"web_socket/src/online_medicine_store/models"
)

type RoomMgr struct {
	Rooms    map[uint32]*Room
	RoomLock sync.RWMutex
}

// RoomMgrObj 全局的房间管理句柄
var RoomMgrObj *RoomMgr

func init() {
	RoomMgrObj = &RoomMgr{
		Rooms: make(map[uint32]*Room),
	}
	models.RDB.Del(define.REDIS_ROOM_COUNT)
	models.RDB.Set(define.REDIS_ROOM_COUNT, 0, 0)
	models.RDB.Del(define.REDIS_ROOM_TALK_USER_NUMBER)
	models.RDB.Set(define.REDIS_ROOM_TALK_USER_NUMBER, 0, 0)
}

func (rm *RoomMgr) AddRoom(room *Room) {
	rm.RoomLock.Lock()
	defer rm.RoomLock.Unlock()
	rm.Rooms[room.ID] = room
}
func (rm *RoomMgr) RemoveRoom(roomid uint32) error {
	rm.RoomLock.Lock()
	err := rm.Rooms[roomid].MasterLeaved()
	if err != nil {
		return err
	}
	fmt.Println(roomid, "号房间被删除了")
	models.RDB.Decr(define.REDIS_ROOM_COUNT)
	delete(rm.Rooms, roomid)
	defer rm.RoomLock.Unlock()
	return nil
}

func (rm *RoomMgr) GetRoomById(roomid uint32) *Room {
	return rm.Rooms[roomid]
}

func (rm *RoomMgr) GetOneRoom() *Room {
	un, err := models.RDB.Get(define.REDIS_ROOM_TALK_USER_NUMBER).Int()
	if err != nil {
		return nil
	}
	rn, err := models.RDB.Get(define.REDIS_ROOM_COUNT).Int()
	if err != nil {
		return nil
	}
	if rn < 1 {
		return nil
	}
	fmt.Println("----------", un, rn)
	roomId := un % rn
	return rm.Rooms[uint32(roomId)]
}

func (rm *RoomMgr) SendMsgToEveryOn(data []byte) {
	for _, v := range rm.Rooms {
		_ = v.SendMsgToAll(data)
	}
}
