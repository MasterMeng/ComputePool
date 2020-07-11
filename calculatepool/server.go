package calculatepool

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	proto "github.com/golang/protobuf/proto"
)

// Server define the server struct
type Server struct {
	port   int
	ws     *restful.WebService
	hp     []*RegisterRequest
	client *http.Client
	hard   int
}

// NewServer return an instance of the server
func NewServer(port int) *Server {
	var hp []*RegisterRequest
	return &Server{
		port:   port,
		ws:     new(restful.WebService),
		hp:     hp,
		client: new(http.Client),
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

// SetHard set the difficulty
func (s *Server) SetHard(hard int) {
	s.hard = hard
}

// GetHard get the difficilty
func (s *Server) GetHard() int {
	return s.hard
}

// Hello for testing if the service is alive
func (s *Server) Hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "hello")
}

// Register register the compute node
func (s *Server) Register(req *restful.Request, resp *restful.Response) {
	body, _ := ioutil.ReadAll(req.Request.Body)

	registerReq := &RegisterRequest{}

	proto.Unmarshal(body, registerReq)

	fmt.Println(registerReq)

	hp := &RegisterRequest{
		Host: registerReq.Host,
		Port: registerReq.Port,
	}
	s.hp = append(s.hp, hp)

	registerResp := &RegisterResponse{
		Info: "Success",
	}

	respBody, _ := proto.Marshal(registerResp)

	resp.ResponseWriter.Write(respBody)
}

// PoW proof of work
func (s *Server) PoW(req *restful.Request, resp *restful.Response) {
	body, _ := ioutil.ReadAll(req.Request.Body)
	powReq := &PoWRequest{}

	proto.Unmarshal(body, powReq)
	s.SetHard(int(powReq.Hard))

	var i int64
	ch := make(chan *PoWResponse)
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	for j := 0; j < cpus; j++ {
		go func() {
			for i = 0; ; i++ {
				resp, hashed := s.calculateHash(powReq.Msg, i)
				fmt.Printf("%s hashed: %s\n", powReq.Msg, hashed)
				if s.isHashValie(hashed) {
					ch <- resp
					break
				}
			}
		}()
	}
	powResp := <-ch

	respBody, _ := proto.Marshal(powResp)
	resp.Write(respBody)
}

// alive get the compute node status
func (s *Server) alive(hp *RegisterRequest) bool {
	url := "http://" + hp.Host + ":" + strconv.Itoa(int(hp.Port)) + "/hello"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := s.client.Do(req)

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// calculateHash calculate the hash
func (s *Server) calculateHash(msg string, i int64) (*PoWResponse, string) {
	now := time.Now().String()
	record := msg + strconv.Itoa(int(i)) + now
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return &PoWResponse{
		Msg:    msg,
		Number: i,
		Time:   now,
	}, hex.EncodeToString(hashed)
}

// isHashValie judge the hash string
func (s *Server) isHashValie(hash string) bool {
	i := len(hash)
	var j int
	for j = 0; j <= i; j++ {
		if hash[j] != '0' {
			break
		}
	}
	return j > s.hard
}

// DoWork assigns PoW to registered node
func (s *Server) DoWork(req *restful.Request, resp *restful.Response) {
	ch := make(chan *http.Response)
	for _, hp := range s.hp {
		go func(hp *RegisterRequest) {
			if s.alive(hp) {
				url := "http://" + hp.Host + ":" + strconv.Itoa(int(hp.Port)) + "/pow"
				powReq, _ := http.NewRequest(http.MethodGet, url, req.Request.Body)
				powReq.Header.Set("Content-Type", "application/json")
				powResp, _ := s.client.Do(powReq)
				ch <- powResp
			}
		}(hp)
	}
	powResp := <-ch

	respBody, _ := ioutil.ReadAll(powResp.Body)
	resp.Write(respBody)
}
