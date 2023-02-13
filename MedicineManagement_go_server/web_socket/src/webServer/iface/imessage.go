package iface

// IMessage 将请求的消息封装到一个message中，定义一个抽象的接口
type IMessage interface {
	GetMsgId() uint32 // 获取消息的ID
	GetData() []byte
	SetMsgId(uint32)
	SetData([]byte)
}
