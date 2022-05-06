package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx_learnin/zinx_learning/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex
}

//创建当前连接的方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}

}

// Add 添加连接
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID=", conn.GetConnID(), "connection add to ConnManager successfully:conn num=", c.Len())
}

// Remove 删除连接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("ConnID=", conn.GetConnID(), "connection remove from ConnManager successfully:conn num=", c.Len())

}

// Get 根据connID获取连接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源Map 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}

}

// Len 得到当前连接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// ClearConn 清除并终止所有连接
func (c *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}
