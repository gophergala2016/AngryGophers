package engine

import "sync/atomic"

const canvasSizeX float32 = 800
const canvasSizeY float32 = 800
const canvasMapX float32 = 16
const canvasMapY float32 = 16

func (s *Server) CalcAll(isClient bool) {
	tmp := refreshModifier
	if isClient {
		refreshModifier = ClientToServerRatio
	}
	s.checkBulletsOnMap(canvasSizeX, canvasSizeY, refreshModifier)

forLoop:
	for _, c := range s.clients {
		if c.GetDeath() {
			continue forLoop
		}
		hit, hitClientId := s.checkHitTank(c)
		if hit {
			life := atomic.AddInt32(&c.Life, -5)
			if life <= 0 {
				atomic.StoreInt32(&c.Life, 100)
				c.SetDeath(true, 0, 0)
				s.explosionAdd(c.PositionX, c.PositionY)
				s.scoreAdd(hitClientId)
				// s.sendResponse("MAP", c.RemoteAddr, s.BuildAnswer(c.id, false))
				continue forLoop
			} else {
				s.newHit = true
			}
		}

		powerUpBool, powerUpId := s.checkPowerup(c.id, c.PositionX, c.PositionY, c.PositionX+tankWidth, c.PositionY+tankHeight)
		if powerUpBool {
			if powerUpId == 1 {
				c.Powerup = 200 * int(refreshModifier)
			}
		}
		if c.Powerup > 0 {
			c.Speed = defaultTankSpeed * 5
			c.Powerup--
			if c.Powerup == 0 {
				c.Speed = defaultTankSpeed
			}
		}

		var speed float32
		// var speed = c.Speed * refreshModifier
		if c.Moving {
			speed = s.setSpeedTank(c, refreshModifier)
			//			log.Println(speed)
			newPositionX := c.PositionX
			newPositionY := c.PositionY
			if s.checkColision(c, newPositionX, newPositionY) == false {
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
		}

		if c.Fire {
			if c.LastFire == 0 {
				c.LastFire = 20 * int(refreshModifier)
				s.newBullet = true
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
		if c.Smoke > 0 && c.Smoke <= 150 {
			if c.Smoke >= 100 && (c.Smoke%10) == 0 {
				s.smokeAdd(c.PositionX, c.PositionY)
			}
			c.Smoke--
		}

	}
	refreshModifier = tmp
}
