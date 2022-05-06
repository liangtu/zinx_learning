package ziface

type IRouter interface {
	// PreHandle 在处理conn业务之前的方法Hook
	PreHandle(request IRequest)

	// Handle 在处理conn业务主方法
	Handle(request IRequest)
	
	// PostHandle 在处理业务之后的方法
	PostHandle(request IRequest)
}
