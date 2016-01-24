package engine

import (
	"math/rand"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

const channelBufSize int = 100

var maxId int = 0
var fullLife int32 = 100
var defaultDirection int = 0
var prev string

var firstPosition [][]float32 = [][]float32{
	[]float32{50, 50},
	[]float32{canvasSizeX - 100, canvasSizeY - 100},
	[]float32{50, canvasSizeY - 100},
	[]float32{canvasSizeX - 100, 50}}

type Client struct {
	id            int
	nick          string
	idReqMax      int32
	RemoteAddr    *net.UDPAddr
	RemoteAddrStr string
	server        *Server
	PositionX     float32
	PositionY     float32
	Life          int32
	Death         bool
	Direction     int
	Speed         float32
	Moving        bool
	Fire          bool
	Smoke         int
	LastFire      int
	Powerup       int
}

func (server *Server) NewClient(remoteAddr *net.UDPAddr, nick string, reqId string) (*Client, int) {
	if remoteAddr == nil {
		panic("remoteAddr cannot be nil")
	}
	tmp := server.clients[remoteAddr.String()] // users[clientId]
	if tmp != nil {
		server.sendResponse("LOGIN", remoteAddr, strconv.Itoa(tmp.GetId())+";"+reqId)
		server.sendPastMessages(tmp)
		return nil, tmp.GetId()
	}

	maxId = int(atomic.AddInt32(&server.userId, 1))
	position := firstPosition[maxId%4]
	server.changesServer = true
	server.scoreNewClient(maxId)
	return &Client{
		maxId,
		nick,
		0,
		remoteAddr,
		remoteAddr.String(),
		server,
		float32(position[0]),
		float32(position[1]),
		fullLife,
		false,
		defaultDirection,
		defaultTankSpeed,
		false,
		false,
		0,
		0,
		0,
	}, 0
}

func (c *Client) GetId() int {
	return c.id
}

func (c *Client) SetId(id int) {
	c.id = id
}

func (c *Client) GetNick() string {
	return c.nick
}

func (c *Client) SetNick(nick string) {
	c.nick = nick
}

func (c *Client) SetDeath(status bool, mapSizeX, mapSizeY float32) {
	if status {
		c.Death = true
	} else {
		if c.Death {
			c.Death = false
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			c.PositionX = float32(r.Intn(int(mapSizeX - tankWidth)))
			c.PositionY = float32(r.Intn(int(mapSizeY - tankHeight)))
		}
	}
}

func (c *Client) GetDeath() bool {
	return c.Death
}
