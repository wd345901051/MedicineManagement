package iface

// IRouter 路由接口，路由里的数据都是IRequest
type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法
	PreHandle(req IRequest)
	// Handle 处理conn业务的方法
	Handle(req IRequest)
	// PostHandle 在处理conn业务之后的钩子方法
	PostHandle(req IRequest)
}
