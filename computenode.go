package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/mastermeng/calculatepool/calculatepool"
)

func main() {
	url := "http://localhost:8888/register"

	registerReq := &calculatepool.RegisterRequest{
		Host: "localhost",
		Port: 9999,
	}

	reqBody, _ := proto.Marshal(registerReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, _ := client.Do(req)

	fmt.Println(resp)

	server := calculatepool.NewServer(9999)
	server.RegisterRoutes("/hello", "get", server.Hello)
	server.RegisterRoutes("/pow", "get", server.PoW)
	server.SetHard(2)
	server.Start()
}
