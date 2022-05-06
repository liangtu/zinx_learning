package ziface

type IServer interface {
	Start()
	Stop()
	// Serve 运行
	Serve()
	// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端处理使用
	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnManager
	// SetOnConnStart 注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(connection IConnection))

	// SetOnConnStop 注册OnConnStop 钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStart 调用OnConnStart 钩子函数的方法
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用OnConnStop 钩子函数的方法
	CallOnConnStop(connection IConnection)
}
