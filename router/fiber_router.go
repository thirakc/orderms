package router

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/pallat/micro/order"
)

type FiberRouter struct {
	*fiber.App
}

func NewFiberRouter() Router {
	return &FiberRouter{fiber.New()}
}

func (r *FiberRouter) GET(relativePath string, handler HandlerFunc) {
	r.App.Get(relativePath, func(c *fiber.Ctx) error {
		handler(&FiberContext{c})
		return nil
	})
}

func (r *FiberRouter) POST(relativePath string, handler HandlerFunc) {
	r.App.Post(relativePath, func(c *fiber.Ctx) error {
		handler(&FiberContext{c})
		return nil
	})
}

func (r *FiberRouter) ListenAndServe() func() {

	go func() {
		if err := r.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		stop()
		fmt.Println("shutting down gracefully, press Ctrl+C again to force")

		if err := r.Shutdown(); err != nil {
			fmt.Println(err)
		}
	}
}

type FiberContext struct {
	*fiber.Ctx
}

func (c *FiberContext) Order() (o order.Order, err error) {
	err = c.Ctx.BodyParser(o)
	return
}

func (c *FiberContext) JSON(code int, v interface{}) {
	c.Ctx.Status(code).JSON(v)
}

func (c *FiberContext) Status(code int) {
	c.Ctx.Status(code)
}
