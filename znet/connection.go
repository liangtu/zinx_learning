package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx_learnin/zinx_learning/utils"
	"zinx_learnin/zinx_learning/ziface"
)

type Connection struct {
	//当前conn属于哪个server
	TcpServer ziface.IServer

	//当前连接的socket套接字
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经退出 channel
	ExitChan chan bool

	//无缓冲管道，用于读写之间的消息通信
	msgChan chan []byte

	//该链接处理的方法
	MsgHandler ziface.IMsgHandle

	//连接属性集合
	property map[string]interface{}

	//保护连接属性的锁
	propertyLock sync.RWMutex
}

//NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandle,
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartReader 连接读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...ConnID =", c.ConnID)
	fmt.Println("--------------------------------------------------")
	defer fmt.Println("connID ", c.ConnID, "Reader is exit; remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitChan <- true
			return
		}

		//拆包 得到msgID 和 msgDataLen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitChan <- true
			return
		}

		//根据datalen再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitChan <- true
				return
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//是否开启了工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {

			//从路由中找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID找到对应处理api的业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写消息
func (c *Connection) StartWriter() {
	fmt.Println("StartWriter is running ...ConnID = ", c.ConnID)
	fmt.Println("--------------------------------------------------")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]！！！")
	//不断的阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan: //代表Reader已经退出，此时Writer也要退出
			return

		}

	}
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)
	//启动从当前连接读数据
	go c.StartReader()
	//从管道读数据
	go c.StartWriter()

	//创建连接之后需要调用的处理业务的钩子方法
	c.TcpServer.CallOnConnStart(c)

}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写给管道
	c.msgChan <- msg
	return nil
}

// Stop 停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID=", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用开发者注册的 销毁连接之前的 需要执行业务的hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitChan <- true
	//将当前连接从ConnMgr中去掉
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭连接管道
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接的绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端 TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
