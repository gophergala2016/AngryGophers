package engine

const canvasSizeX float32 = 800
const canvasSizeY float32 = 800
const canvasMapX float32 = 16
const canvasMapY float32 = 16

func (s *Server) calcAll() {
	s.checkBulletsOnMap(canvasSizeX, canvasSizeY, refreshModifier)

forLoop:
	for _, c := range s.clients {
		hit, hitClientId := s.checkHitTank(c)
		if hit {
			s.explosionAdd(c.PositionX, c.PositionY)
			c.PositionX = c.StartPosX
			c.PositionY = c.StartPosY
			s.scoreAdd(hitClientId)
			s.sendResponse("MAP", c.RemoteAddr, s.BuildAnswer(c.id, false))
			continue forLoop
		}

		var speed float32
		// var speed = c.Speed * refreshModifier
		if c.Moving {
			speed = s.setSpeedTank(c, refreshModifier)
			newPositionX := c.PositionX
			newPositionY := c.PositionY
			switch c.Direction {
			case 0:
				newPositionY = c.PositionY - speed
				if newPositionY <= 0 {
					newPositionY = 0
				}
			case 90:
				newPositionX = c.PositionX + speed
				if newPositionX+tankHeight >= canvasSizeX {
					newPositionX = canvasSizeX - tankHeight
				}
			case 180:
				newPositionY = c.PositionY + speed
				if newPositionY+tankHeight >= canvasSizeY {
					newPositionY = canvasSizeY - tankHeight
				}
			case 270:
				newPositionX = c.PositionX - speed
				if newPositionX <= 0 {
					newPositionX = 0
				}
			}
			s.checkColision(c, newPositionX, newPositionY)
		}

		if c.Fire {
			if c.LastFire == 0 {
				c.LastFire = 20 * int(refreshModifier)
				s.bullets = append(s.bullets,
					&Bullet{
						speed:     speed + bulletSpeed,
						ownerId:   c.id,
						x:         c.PositionX + tankWidthHalf - bulletWidthHalf,
						y:         c.PositionY + tankHeightHalf - bulletHeightHalf,
						direction: c.Direction})
			}
		}
		if c.LastFire > 0 {
			c.LastFire--
		}

	}
}
