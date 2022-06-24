package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pallat/micro/order"
)

type GinRouter struct {
	*gin.Engine
}

func New() Router {
	return &GinRouter{gin.Default()}
}

func (r *GinRouter) GET(relativePath string, handler HandlerFunc) {
	r.Engine.GET(relativePath, func(c *gin.Context) {
		handler(&GinContext{c})
	})
}

func (r *GinRouter) POST(relativePath string, handler HandlerFunc) {
	r.Engine.POST(relativePath, func(c *gin.Context) {
		handler(&GinContext{c})
	})
}

func (r *GinRouter) ListenAndServe() func() {
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		stop()
		fmt.Println("shutting down gracefully, press Ctrl+C again to force")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(timeoutCtx); err != nil {
			fmt.Println(err)
		}
	}
}

type GinContext struct {
	*gin.Context
}

func (c *GinContext) Order() (o order.Order, err error) {
	err = c.ShouldBindJSON(&o)
	return
}

func (c *GinContext) JSON(code int, v interface{}) {
	c.Context.JSON(code, v)
}

func (c *GinContext) Status(code int) {
	c.Context.Status(code)
}
