package net

import (
	"fmt"
	"web_socket/src/webServer/iface"
	"web_socket/src/webServer/utils"
)

type MsgHandle struct {
	// 存放每一个MsgID 所对应的处理方法
	Apis map[uint32]iface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中读取
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 执行对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(req iface.IRequest) {
	// 从req中找到msgID
	handler, ok := mh.Apis[req.GetMsg().GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", req.GetMsg().GetMsgId(), "is NOT FOUND! Need Register")
		return
	}
	// 根据msgID调度对应的router业务
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router iface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("repeat api , msgID = ", msgID)
		return
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, "success!")
}

// StartWorkerPool 启动一个Worker工作池(只能发生一次)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据配置分别开启worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 当前的Worker对应的channel消息队列开辟空间,地0个就用第0个
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 启动当前的worker，阻塞等待消息从channel中传递
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, TaskQueue chan iface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")
	for {
		select {
		// 如果由消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case req := <-TaskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(req iface.IRequest) {
	// 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := req.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", req.GetConnection().GetConnID(),
		" req msgID = ", req.GetMsg().GetMsgId(), "to workerID = ", workerID)
	// 将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- req
}
