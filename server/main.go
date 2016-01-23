package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"../engine"
)

const framePerSec int64 = 30
const timePerFrame int64 = 1000000 / framePerSec

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	_, err := conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	ser, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	log.Println(timePerFrame, " - timePerFrame")

	server := engine.NewServer(ser)
	go server.Listen()
	log.Println("listening on 8081...")
	go func() {
		for {
			actualTime := time.Now().UnixNano()

			server.SendAll()
			log.Print("timeNow", time.Now())

			differenceTime := (time.Now().UnixNano() - actualTime) / 1000 //microseconds
			//log.Print(differenceTime)
			if differenceTime < timePerFrame {
				//	log.Println("Sleeep", int64((timePerFrame-differenceTime)/1000))
				//	log.Println(time.Duration(timePerFrame-differenceTime) * time.Microsecond)
				time.Sleep(time.Duration(timePerFrame-differenceTime) * time.Microsecond)
			}
		}

	}()

	for {
		msg := make([]byte, 100)
		n, remoteaddr, err := ser.ReadFromUDP(msg)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, msg)

		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		tmp := strings.Split(string(msg[:n]), ";")
		if len(tmp) > 1 {
			switch string(tmp[1]) {
			case "login":
				if len(tmp) == 3 {
					client, clientId := server.NewClient(remoteaddr, tmp[2], tmp[0])
					if clientId == 0 {
						server.Add(client, tmp[0])
					}
				}
			default:
				server.ParseResponse(tmp[0], tmp[1], remoteaddr)
			}
		}
	}

}
