package ziface

type name interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error) //打包
	Unpack([]byte) (IMessage, error)   //拆包
}
