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
	// for i := 0; i < 5; i++ {
	x := float32(r.Intn(int(mapSizeX - tankWidth)))
	y := float32(r.Intn(int(mapSizeY - tankHeight)))
	s.powerups.powerup = append(s.powerups.powerup,
		&Powerup{
			x:   x,
			y:   y,
			typ: 1,
		})
	// log.Println("X Y ", x, y)
	// }

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
			s.generatePowerup(canvasSizeX, canvasSizeY)
			s.powerups.show = true
			return true, b.typ
		}
	}
	return false, 0
}
