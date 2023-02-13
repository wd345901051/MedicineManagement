package net

import "web_socket/src/webServer/iface"

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需求需要对这个基类的方法进行重写就好了
type BaseRouter struct{}

// PreHandle 在处理conn业务之前的钩子方法
func (br *BaseRouter) PreHandle(req iface.IRequest) {}

// Handle 处理conn业务的方法
func (br *BaseRouter) Handle(req iface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法
func (br *BaseRouter) PostHandle(req iface.IRequest) {}
