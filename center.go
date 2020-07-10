package main

import (
	"fmt"

	"github.com/mastermeng/computepool/computepool"
)

func main() {
	fmt.Println("hello")
	server := computepool.NewServer(8888)
	server.RegisterRoutes("/hello", "get", server.Hello)
	server.RegisterRoutes("/register", "post", server.Register)
	server.RegisterRoutes("/dowork", "get", server.DoWork)
	server.Start()
}
