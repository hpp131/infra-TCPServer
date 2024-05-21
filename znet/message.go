package znet

func NewMessage(data []byte, id uint32) *Message {
	return &Message{
		Data: data,
		MsgID: id,
		Len: uint32(len(data)),
	}
}

type Message struct {
	Data []byte
	MsgID uint32
	Len uint32
}

// Implement IMessage interface

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetDataLen() uint32 {
	return m.Len
}

func (m *Message) SetData(data []byte) ()  {
	m.Data = data
}
func (m *Message) SetMsgID(msgID uint32) ()  {
	m.MsgID = msgID
}

func (m *Message) SetDataLen(length uint32)  {
	m.Len = length
}