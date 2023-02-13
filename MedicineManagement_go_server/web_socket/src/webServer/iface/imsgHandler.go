package iface

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	// DoMsgHandler 执行对应的router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
	SendMsgToTaskQueue(req IRequest)
}
