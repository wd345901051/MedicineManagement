package net

type Message struct {
	Id   uint32 `json:"id"`   //消息的ID
	Data []byte `json:"data"` //消息的内容
}

// NewMsgPackage 创建一个Msg的方法
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:   id,
		Data: data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
} // 获取消息的ID

func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
