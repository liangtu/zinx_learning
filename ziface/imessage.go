package ziface

type IMessage interface {
	GetMsgId() uint32
	GetData() []byte
	GetDataLen() uint32

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
