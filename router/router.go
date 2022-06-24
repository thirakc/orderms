package router

import "github.com/pallat/micro/order"

type Router interface {
	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	ListenAndServe() func()
}

type HandlerFunc func(order.Context)
