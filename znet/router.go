package znet

import "zinx_learnin/zinx_learning/ziface"

// BaseRouter 在业务中继承这个基类 去重写你需要的方法 并不需要把全部方法实现 /*
type BaseRouter struct {
}

// PreHandle 在处理conn业务之前的方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

// Handle 在处理conn业务主方法
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

// PostHandle 在处理业务之后的方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
