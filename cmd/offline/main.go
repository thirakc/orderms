package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pallat/micro/order"
	"github.com/pallat/micro/router"
	"github.com/pallat/micro/store"
)

func init() {
	err := godotenv.Load("offline.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}
}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	handler := order.NewHandler(store.NewMariaDBStore(os.Getenv("DSN")), os.Getenv("FILTER_CHANNEL"))

	var r router.Router

	if os.Getenv("ROUTER") == "gin" {
		r = router.New()
	} else {
		r = router.NewFiberRouter()
	}

	r.POST("api/v1/orders", handler.Order)
	r.ListenAndServe()()
}
