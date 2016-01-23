package engine

import (
	"net"
	"strconv"
	"sync/atomic"
)

const channelBufSize int = 100

var maxId int = 0
var fullLife int = 100
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
	Life          int
	Direction     int
	Speed         float32
	Moving        bool
	Fire          bool
	LastFire      int
	StartPosX     float32
	StartPosY     float32
}

func (server *Server) NewClient(remoteAddr *net.UDPAddr, nick string, reqId string) (*Client, int) {
	if remoteAddr == nil {
		panic("remoteAddr cannot be nil")
	}
	tmp := server.clients[remoteAddr.String()] // users[clientId]
	if tmp != nil {
		server.sendResponse("LOGIN", remoteAddr, strconv.Itoa(tmp.GetId())+";"+reqId)
		return nil, tmp.GetId()
	}

	maxId = int(atomic.AddInt32(&server.userId, 1))
	position := firstPosition[maxId%4]

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
		defaultDirection,
		defaultTankSpeed,
		false,
		false,
		0,
		float32(position[0]),
		float32(position[1])}, 0
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
