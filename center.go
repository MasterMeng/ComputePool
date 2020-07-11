package main

import (
	"fmt"

	"github.com/mastermeng/calculatepool/calculatepool"
)

func main() {
	fmt.Println("hello")
	server := calculatepool.NewServer(8888)
	server.RegisterRoutes("/hello", "get", server.Hello)
	server.RegisterRoutes("/register", "post", server.Register)
	server.RegisterRoutes("/dowork", "get", server.DoWork)
	server.Start()
}
