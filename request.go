package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/mastermeng/computepool/computepool"
)

func main() {
	url := "http://localhost:8888/dowork"

	registerReq := &computepool.PowRequest{
		Msg: "localhost",
	}

	reqBody, _ := proto.Marshal(registerReq)

	req, _ := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, _ := client.Do(req)

	fmt.Println(resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)

	pow := &computepool.PowResponse{}
	proto.Unmarshal(respBody, pow)

	fmt.Println(pow)
}
