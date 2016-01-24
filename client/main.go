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
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/net/websocket"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr
var serverMsg chan []byte
var requestCounter int32 = 0
var waitingRequests map[int32][]byte
var clientId int
var clientNick string
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

func receiveServerMessages(ws *websocket.Conn, closeWs chan bool) {
	for {
		select {
		case <-closeWs:
			return
		default:
		}
		serverMessage := <-serverMsg
		serverMessageString := strings.SplitN(string(serverMessage), ";", 3)
		if len(serverMessageString) == 3 {
			serverMessageId, err := strconv.ParseInt(serverMessageString[0], 10, 64)
			CheckError(err)
			if serverMessageId < lastReceivedId {
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
					clientId = id
					clientNick = waitingRequestsArray[1]
					delete(waitingRequests, int32(key))
				}
			case "OK":
				key, err := strconv.ParseInt(serverMessageString[2], 10, 32)
				CheckError(err)
				if waitingRequests[int32(key)] != nil {
					delete(waitingRequests, int32(key))
				}
			case "F":
				//				if clientId > 0 {
				//	log.Println(clientId)
				sendMap := server.ParseMsgFromServerToStruct(serverMessageString[2], clientId)

				if sendMap {
					_, err := ws.Write([]byte(serverMessageString[2]))
					CheckError(err)
				}
				//				}

			}
		}

	}
}

func manageWebSocket(ws *websocket.Conn, closeWs chan bool) {
	for {
		select {
		case <-closeWs:
			return
		default:
		}
		if clientId > 0 {
			actualTime := time.Now().UnixNano()
			server.CalcAll(true)
			server.SendAllClient(clientId, ws)
			//			log.Print("timeNow", time.Now())

			differenceTime := (time.Now().UnixNano() - actualTime) / 1000 //microseconds
			//log.Print(differenceTime)
			if differenceTime < engine.ClientTimePerFrame {
				//	log.Println("Sleeep", int64((timePerFrame-differenceTime)/1000))
				//	log.Println(time.Duration(timePerFrame-differenceTime) * time.Microsecond)
				time.Sleep(time.Duration(engine.ClientTimePerFrame-differenceTime) * time.Microsecond)
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
	server = engine.NewServer(nil)
	log.Println(server)
	serverMsg = make(chan []byte)
	waitingRequests = make(map[int32][]byte)

	ReadFromWebsocket := func(ws *websocket.Conn) {
		closeWs := make(chan bool)
		closeWs2 := make(chan bool)
		go receiveServerMessages(ws, closeWs)
		go manageWebSocket(ws, closeWs2)
	forLoop:
		for {
			msg := make([]byte, 100)
			n, err := ws.Read(msg)

			if err != nil {
				closeWs <- true
				closeWs2 <- true
				break forLoop
			}
			messageToSend := msg[:n]
			if string(messageToSend) == "check" && clientNick == "" {
				continue forLoop
			}

			msgId := atomic.AddInt32(&requestCounter, 1)

			go func(messageToSend []byte, msgId int32) {
				if string(messageToSend) == "check" && clientNick != "" {
					messageToSend = []byte("login;" + clientNick)
				}

				if clientId > 0 {
					tmp := server.SelectClient(clientId)
					if tmp != nil {
						engine.UpdateClientData(tmp, string(messageToSend))
					}
				}

				waitingRequests[msgId] = messageToSend
				for waitingRequests[msgId] != nil {
					toSend := []byte(fmt.Sprintf("%d;", msgId))
					sendMessage(append(toSend, messageToSend...))
					<-time.After(time.Second)
				}
			}(messageToSend, msgId)
		}
	}
	http.Handle("/echo", websocket.Handler(ReadFromWebsocket))

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	ServerAddr, err := net.ResolveUDPAddr("udp", "40.76.43.250:8081")
	CheckError(err)

	ClientAddr, err := net.ResolveUDPAddr("udp", ":0")
	CheckError(err)

	conn, err = net.DialUDP("udp", ClientAddr, ServerAddr)
	CheckError(err)
	defer conn.Close()
	go manageMessages()

	open.Run("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
