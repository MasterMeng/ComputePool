package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/mastermeng/computepool/computepool"
)

func main() {
	url := "http://localhost:8888/register"

	registerReq := &computepool.RegisterRequest{
		Host: "localhost",
		Port: 9999,
	}

	reqBody, _ := proto.Marshal(registerReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, _ := client.Do(req)

	fmt.Println(resp)

	server := computepool.NewServer(9999)
	server.RegisterRoutes("/hello", "get", server.Hello)
	server.RegisterRoutes("/pow", "get", server.PoW)
	server.SetHard(2)
	server.Start()
}
