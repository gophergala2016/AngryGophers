package engine

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
)

func (s *Server) ParseResponse(idReq string, msg string, remoteaddr *net.UDPAddr) {
	tmp := s.clients[remoteaddr.String()] // users[clientId]
	if tmp == nil {
		log.Print("no user found - ", remoteaddr.String())
		s.sendResponse("ERROR", remoteaddr, "no user found")
		return
	}

	idReqInt, err := strconv.ParseInt(idReq, 10, 32)
	if err != nil {
		log.Println(err)
		return
	}
	if tmp.idReqMax < int32(idReqInt) {
		tmp.idReqMax = int32(idReqInt)
	} else {
		log.Print("message is old")
		return
	}
	switch msg {
	case "fire":
		tmp.Fire = true
	case "fire2":
		tmp.Fire = false
	case "right":
		tmp.Direction = 90
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "left":
		tmp.Direction = 270
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "down":
		tmp.Direction = 180
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "up":
		tmp.Direction = 0
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "right2", "left2", "down2", "up2":
		tmp.Moving = false
		tmp.Speed = 0
	}

	s.clients[remoteaddr.String()] = tmp
	s.sendResponse("OK", remoteaddr, idReq)
}

func (self *Server) BuildAnswer(clientId int, firstAnswer bool) string {
	var result bytes.Buffer
	if firstAnswer {
		x := self.mapa.drawMap()
		for _, v := range x {
			result.WriteString("M;")
			for _, v2 := range v {
				result.WriteString(strconv.Itoa(v2) + ";")
			}
			result.WriteString("\n")
		}
	}
	for _, u := range self.bullets {
		result.WriteString(fmt.Sprintf("B;%.0f;%.0f;%d;\n",
			u.x, u.y, u.direction))
	}
	for _, user := range self.clients {
		color := "r"
		if clientId == user.id {
			color = "b"
		}

		result.WriteString(fmt.Sprintf("T;%d;%s;%.0f;%.0f;%.0f;%d;%d;%d;\n",
			user.id, color, user.PositionX, user.PositionY, user.Speed, user.Direction, user.Direction, 100))
	}
	if self.explosion.show {
		for _, point := range self.explosion.position {
			result.WriteString(fmt.Sprintf("E;%.0f;%.0f;\n", point.x, point.y))
		}
	}
	if self.score.change {
		for id, point := range self.score.client {
			result.WriteString(fmt.Sprintf("S;%d;%d;\n", id, point))
		}
	}
	return result.String()
}

/*
Odpowiedz format
tank
obiekt;id;color;pozycjaX;pozycjaY;obrot;obrot_lufy;zycie(hp);
T;1;R;10;10;0;0;50;

kolor R G B K

*/
