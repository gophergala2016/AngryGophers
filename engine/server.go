package engine

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

var refreshModifier float32 = 1

type Server struct {
	conn      *net.UDPConn
	userId    int32
	reqId     int32
	clients   map[string]*Client
	bullets   []*Bullet
	explosion Explosion
	addCh     chan *Client
	doneCh    chan bool
	errCh     chan error
	score     Scores
	mapa      *mapa
}

func NewServer(conn *net.UDPConn) *Server {
	var bullets []*Bullet
	explosion := Explosion{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	doneCh := make(chan bool)
	errCh := make(chan error)
	var score Scores
	score.client = make(map[int]int)


	var speedGround []int = []int{10, 0, 10, 30, 3}
	var mapa1 [][]int = [][]int{
		[]int{2, 2, 2, 0, 0, 2, 0, 1, 2, 0, 0, 2, 0, 2, 0, 2},
		[]int{0, 2, 0, 0, 2, 0, 2, 1, 2, 2, 0, 2, 0, 2, 2, 0},
		[]int{0, 2, 0, 0, 2, 2, 2, 1, 2, 0, 2, 2, 0, 2, 0, 2},
		[]int{0, 2, 0, 0, 2, 0, 2, 1, 2, 0, 0, 2, 0, 2, 0, 2},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
		[]int{0, 2, 0, 0, 0, 2, 0, 1, 2, 0, 0, 2, 0, 2, 2, 2},
		[]int{2, 0, 0, 0, 2, 0, 2, 1, 2, 2, 2, 2, 0, 2, 0, 0},
		[]int{2, 2, 2, 0, 2, 2, 2, 1, 2, 0, 0, 2, 0, 2, 2, 2},
		[]int{0, 2, 2, 0, 2, 0, 2, 1, 2, 0, 0, 2, 0, 2, 2, 2},
	}

	m := getMap(mapa1, speedGround)
	var reqId int32 = 0
	var userId int32 = 1
	s := &Server{
		conn,
		userId,
		reqId,
		clients,
		bullets,
		explosion,
		addCh,
		doneCh,
		errCh,
		score,
		m,
	}

//	s.mapa.setMap()
	return s
}

func (s *Server) Add(c *Client, reqId string) {
	s.sendResponse("LOGIN", c.RemoteAddr, strconv.Itoa(c.GetId())+";"+reqId)
	s.addCh <- c
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	s.sendResponse("F", c.RemoteAddr, s.BuildAnswer(c.id, true))
}

func (s *Server) SendAll() {
	s.calcAll()
	for _, c := range s.clients {
		m := s.BuildAnswer(c.id, false)
		s.sendResponse("F", c.RemoteAddr, m)
	}
	s.scoreRead()
	s.explosionRead()
}

func (s *Server) Listen() {
	for {
		select {
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.RemoteAddrStr] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			log.Println("doneCh")
			return
		}
	}
}

func (s *Server) sendResponse(typ string, addr *net.UDPAddr, msg string) {
	id := atomic.AddInt32(&s.reqId, 1)
	// log.Print("msg: ", fmt.Sprintf("%d;%s;%s", id, typ, msg))
	_, err := s.conn.WriteToUDP([]byte(fmt.Sprintf("%d;%s;%s", id, typ, msg)), addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}
