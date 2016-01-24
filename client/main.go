package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"../engine"
	"golang.org/x/net/websocket"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr
var serverMsg chan []byte
var requestCounter int32 = 0
var waitingRequests map[int32][]byte
var client engine.Client
var server *engine.Server
var currentMap *engine.Mapa = &engine.Mapa{}
var lastReceivedId int64 = 0

//// CLIENT ////////
func sendMessage(msg []byte) {
	log.Println(string(msg))
	_, err := conn.Write(msg)
	CheckError(err)
}

func manageMessages() {
	for {
		msgFromServer := make([]byte, 2048)
		n, err := bufio.NewReader(conn).Read(msgFromServer)
		CheckError(err)
		serverMsg <- msgFromServer[:n]
	}
}

func manageWebSocket(ws *websocket.Conn, closeWs chan bool) {
	for {
		select {
		case <-closeWs:
			return
		default:
		}
		serverMessage := <-serverMsg
		serverMessageString := strings.SplitN(string(serverMessage), ";", 3)
		// log.Println(serverMessageString)
		// log.Println(waitingRequests)
		if len(serverMessageString) == 3 {
			serverMessageId, err := strconv.ParseInt(serverMessageString[0], 10, 64)
			CheckError(err)
			if(serverMessageId < lastReceivedId){
				continue
			}
			lastReceivedId = serverMessageId
			
			switch serverMessageString[1] {
			case "LOGIN":
				idKey := strings.SplitN(string(serverMessageString[2]), ";", 2)
				idString := idKey[0]
				keyString := idKey[1]
				key, err := strconv.ParseInt(keyString, 10, 32)
				log.Println(key)
				CheckError(err)
				if waitingRequests[int32(key)] != nil {
					waitingRequestsArray := strings.Split(string(waitingRequests[int32(key)]), ";")
					id, err := strconv.Atoi(idString)
					CheckError(err)
					client.SetId(id)
					client.SetNick(waitingRequestsArray[1])
					delete(waitingRequests, int32(key))
				}
			case "OK":
				key, err := strconv.ParseInt(serverMessageString[2], 10, 32)
				CheckError(err)
				if waitingRequests[int32(key)] != nil {
					delete(waitingRequests, int32(key))
				}
			case "F":
				//				log.Println(serverMessageString[2])
				lines := strings.Split(serverMessageString[2], "\n")
				for _, line := range lines {

					if len(line) > 0 {
						cols := strings.Split(line, ";")
						switch cols[0] {
						case "M":
							switch cols[1] {
							case "1":
								currentMap = engine.GetMap(engine.Mapa1, engine.SpeedGround1, 800, 800)
							}
						}
					}

				}

				_, err := ws.Write([]byte(serverMessageString[2]))
				CheckError(err)

			}
		}

	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	client = engine.Client{}
	server = engine.NewServer(nil)
	serverMsg = make(chan []byte)
	waitingRequests = make(map[int32][]byte)
	ReadFromWebsocket := func(ws *websocket.Conn) {
		closeWs := make(chan bool)
		go manageWebSocket(ws, closeWs)
	forLoop:
		for {
			msg := make([]byte, 100)
			n, err := ws.Read(msg)

			if err != nil {
				closeWs <- true
				break forLoop
			}
			messageToSend := msg[:n]
			if string(messageToSend) == "check" && client.GetNick() == "" {
				continue forLoop
			}
			
			msgId := atomic.AddInt32(&requestCounter, 1)
			
			go func(messageToSend []byte, msgId int32) {
				if string(messageToSend) == "check" && client.GetNick() != "" {
					messageToSend = []byte("login;" + client.GetNick())
				}
				
				waitingRequests[msgId] = messageToSend
				for waitingRequests[msgId] != nil {
					toSend := []byte(fmt.Sprintf("%d;", msgId))
					sendMessage(append(toSend, messageToSend...))
					<-time.After(time.Second * 2)
				}
			}(messageToSend, msgId)
		}
	}
	http.Handle("/echo", websocket.Handler(ReadFromWebsocket))

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	ServerAddr, err := net.ResolveUDPAddr("udp", "89.72.59.44:8081")
	CheckError(err)

	ClientAddr, err := net.ResolveUDPAddr("udp", ":0")
	CheckError(err)

	conn, err = net.DialUDP("udp", ClientAddr, ServerAddr)
	CheckError(err)
	defer conn.Close()

	go manageMessages()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
