package engine

import (
	"math/rand"
	"time"
)

type Powerups struct {
	show    bool
	powerup []*Powerup
}

type Powerup struct {
	x   float32
	y   float32
	typ int
}

func (s *Server) generatePowerup(mapSizeX, mapSizeY float32) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		s.powerups.powerup = append(s.powerups.powerup,
			&Powerup{
				x:   float32(r.Intn(int(mapSizeX - tankWidth))),
				y:   float32(r.Intn(int(mapSizeY - tankHeight))),
				typ: 1,
			})
	}

	s.powerups.show = true
}

func (s *Server) powerupRead() {
	s.powerups.show = false
}

func (s *Server) checkPowerup(clientId int, tankX1, tankY1, tankX2, tankY2 float32) (bool, int) {
	for k, b := range s.powerups.powerup {
		if (tankX2 > b.x && tankX1 < b.x) && (tankY2 > b.y && tankY1 < b.y) {
			var tmpList []*Powerup
			for k2, b2 := range s.powerups.powerup {
				if k != k2 {
					tmpList = append(tmpList, b2)
				}
			}
			s.powerups.powerup = tmpList
			s.powerups.show = true
			return true, b.typ
		}
	}
	return false, 0
}
