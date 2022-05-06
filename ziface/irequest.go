package ziface

/**
IRequest接口：
把客户端请求连接信息和请求数据包装到了一个Request中
*/

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到请求数据
	GetData() []byte
	// GetMsgID 得到请求的消息ID
	GetMsgID() uint32
}
