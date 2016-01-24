package engine

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
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
	case "respawn":
		tmp.SetDeath(false, canvasSizeX, canvasSizeY)
	case "fire":
		tmp.Fire = true
	case "fire2":
		tmp.Fire = false
	case "smoke":
		if tmp.Smoke == 0 {
			tmp.Smoke = 150
		}
	// case "smoke2":
	// 	tmp.Smoke = false
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
		x, speeds, name := self.mapa.drawMap()
		result.WriteString("MN;" + name + ";\n")
		for _, v := range x {
			result.WriteString("M;")
			for _, v2 := range v {
				result.WriteString(strconv.Itoa(v2) + ";")
			}
			result.WriteString("\n")
		}
		result.WriteString("MS;")
		for _, v := range speeds {
			result.WriteString(strconv.Itoa(v) + ";")
		}
		result.WriteString("\n")

		for _, v := range self.mapa.GetTrees() {
			result.WriteString(fmt.Sprintf("MT;%d;%d;\n", v[0], v[1]))
		}

		for _, v := range self.mapa.GetRocks() {
			result.WriteString(fmt.Sprintf("MR;%d;%d;\n", v[0], v[1]))
		}

	}
	for _, u := range self.bullets {
		result.WriteString(fmt.Sprintf("B;%.0f;%.0f;%d;%.2f;%d;\n",
			u.x, u.y, u.direction, u.speed, u.ownerId))
	}

	for _, user := range self.clients {
		typ := ""
		if user.Death {
			if clientId == user.id {
				result.WriteString("X;\n")

			}
			typ = "X"
		}
		color := "r"
		if clientId == user.id {
			color = "b"
		}

		result.WriteString(fmt.Sprintf("T%s;%d;%s;%.0f;%.0f;%.0f;%d;%d;%d;\n",
			typ, user.GetId(), color, user.PositionX, user.PositionY, user.Speed, user.Direction, user.Direction, user.Life))

	}

	if self.explosion.show || firstAnswer {
		for _, point := range self.explosion.position {
			result.WriteString(fmt.Sprintf("E;%.0f;%.0f;\n", point.x, point.y))
		}
	}
	if self.smoke.show || firstAnswer {
		for _, point := range self.smoke.position {
			result.WriteString(fmt.Sprintf("SMOKE;%.0f;%.0f;\n", point.x, point.y))
		}
	}
	if self.score.change || firstAnswer {
		for id, point := range self.score.client {
			result.WriteString(fmt.Sprintf("S;%d;%d;\n", id, point))
		}
	}
	if self.changesServer || firstAnswer {
		for _, user := range self.clients {
			result.WriteString(fmt.Sprintf("U;%d;%s;\n",
				user.GetId(), user.GetNick()))
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

func (s *Server) ParseMsgFromServerToStruct(msg string, clientId int) {
	var tmpBullets []*Bullet
	var tmpExplosion []*Position
	var tmpSmoke []*Position
	tmpTanks := make(map[string]*Client)
	tmpUserNick := make(map[int]string)
	tmpScore := make(map[int]int)

	lines := strings.Split(msg, "\n")
forline:
	for _, line := range lines {
		data := strings.Split(line, ";")
		switch data[0] {
		case "B":
			x, err := strconv.ParseFloat(data[1], 32)
			if err != nil {
				continue forline
			}
			y, err := strconv.ParseFloat(data[2], 32)
			if err != nil {
				continue forline
			}
			direction, err := strconv.Atoi(data[3])
			if err != nil {
				continue forline
			}
			speed, err := strconv.ParseFloat(data[4], 32)
			if err != nil {
				continue forline
			}
			ownerId, err := strconv.Atoi(data[5])
			if err != nil {
				continue forline
			}
			tmpBullets = append(tmpBullets, &Bullet{x: float32(x), y: float32(y), speed: float32(speed), direction: direction, ownerId: ownerId})

		case "X":
			for k, c := range s.clients {
				if clientId == c.id {
					s.clients[k].Death = true
				}
			}

		case "T", "TX":
			userId, err := strconv.Atoi(data[1])
			if err != nil {
				continue forline
			}

			positionX, err := strconv.ParseFloat(data[3], 32)
			if err != nil {
				continue forline
			}
			positionY, err := strconv.ParseFloat(data[4], 32)
			if err != nil {
				continue forline
			}
			life, err := strconv.Atoi(data[8])
			if err != nil {
				continue forline
			}
			direction, err := strconv.Atoi(data[6])
			if err != nil {
				continue forline
			}
			speed, err := strconv.ParseFloat(data[5], 32)
			if err != nil {
				continue forline
			}
			moving := false
			if speed > 0 {
				moving = true
			}
			death := false
			if data[0] == "TX" {
				death = true
			}
			tmpTanks[data[1]] = &Client{
				id:        userId,
				nick:      data[2],
				PositionX: float32(positionX),
				PositionY: float32(positionY),
				Life:      life,
				Death:     death,
				Direction: direction,
				Speed:     float32(speed),
				Moving:    moving,
			}

		case "E":
			x, err := strconv.ParseFloat(data[1], 32)
			if err != nil {
				continue forline
			}
			y, err := strconv.ParseFloat(data[2], 32)
			if err != nil {
				continue forline
			}

			tmpExplosion = append(tmpExplosion, &Position{x: float32(x), y: float32(y)})

		case "SMOKE":
			x, err := strconv.ParseFloat(data[1], 32)
			if err != nil {
				continue forline
			}
			y, err := strconv.ParseFloat(data[2], 32)
			if err != nil {
				continue forline
			}

			tmpSmoke = append(tmpSmoke, &Position{x: float32(x), y: float32(y)})

		case "S":
			id, err := strconv.Atoi(data[1])
			if err != nil {
				continue forline
			}
			point, err := strconv.Atoi(data[2])
			if err != nil {
				continue forline
			}
			tmpScore[id] = point
			
		case "M":
			id, err := strconv.Atoi(data[1])
			if err != nil {
				continue forline
			}
			switch id {
				case 1:
				s.SetMap(Mapa1, SpeedGround1)
			}
			
		case "U":
			id, err := strconv.Atoi(data[1])
			if err != nil {
				continue forline
			}
			tmpUserNick[id] = data[2]
		}
	}

	if len(tmpExplosion) > 0 {
		s.explosion.position = tmpExplosion
		s.explosion.show = true
	}

	if len(tmpSmoke) > 0 {
		s.smoke.position = tmpSmoke
		s.smoke.show = true
	}

	if len(tmpScore) > 0 {
		s.score.client = tmpScore
		s.score.change = true
	}

	s.bullets = tmpBullets
	s.clients = tmpTanks
	for k, c := range s.clients {
		if tmpUserNick[c.id] != "" {
			s.clients[k].nick = tmpUserNick[c.id]
		}
	}
}
