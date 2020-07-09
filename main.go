package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

// Server define the server struct
type Server struct {
	port int
	ws   *restful.WebService
	hp   map[string]int
}

// NewServer return an instance of the server
func NewServer(port int) *Server {
	return &Server{
		port: port,
		ws:   new(restful.WebService),
		hp:   make(map[string]int),
	}
}

// Start run the server
func (s *Server) Start() {
	restful.Add(s.ws)
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

// RegisterRoutes register routes to server
func (s *Server) RegisterRoutes(route, method string, function func(*restful.Request, *restful.Response)) {
	switch strings.ToLower(method) {
	case "post":
		s.ws.Route(s.ws.POST(route).To(function))
	case "delete":
		s.ws.Route(s.ws.DELETE(route).To(function))
	case "put":
		s.ws.Route(s.ws.PUT(route).To(function))
	case "get":
		s.ws.Route(s.ws.GET(route).To(function))
	}
}

func main() {
	fmt.Println("hello")

	server := NewServer(8888)
	server.RegisterRoutes("/hello", "get", hello)
	server.Start()
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "hello")
}
