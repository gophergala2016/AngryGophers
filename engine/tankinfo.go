package engine

import (
	"math"
)

const tankWidth float32 = 37.5
const tankWidthHalf float32 = 18.75
const tankHeight float32 = 35
const tankHeightHalf float32 = 17.5

const defaultTankSpeed float32 = 5

func (s *Server) checkHitTank(c *Client) (bool, int) {
	hit, hitClientId := s.checkHitBullet(c.id, c.PositionX-10, c.PositionY-10, c.PositionX+tankWidth, c.PositionY+tankHeight)
	if hit {
		return true, hitClientId
	}
	return false, 0
}

func (s *Server) setSpeedTank(c *Client, refreshTime float32) float32 {
	return c.Speed * refreshTime * s.getSpeedPosition(c.PositionX+tankWidthHalf, c.PositionY+tankHeightHalf)
}

func (s *Server) checkColision(c *Client, newX, newY float32) bool {
	if s.getCollision(newX+tankWidthHalf, newY+tankHeightHalf) {
		return true
	}
	if s.checkAnotherTankColision(c.GetId(), newX+tankWidthHalf, newY+tankHeightHalf) {
		return true
	}

	c.PositionX = newX
	c.PositionY = newY
	return false
}

func (s *Server) checkAnotherTankColision(cliendId int, x1, y1 float32) bool {
	for _, c := range s.clients {
		if c.id != cliendId && c.Death == false {
			x2 := c.PositionX + tankWidthHalf
			y2 := c.PositionY + tankHeightHalf

			if float32(math.Abs(float64(x2-x1))) < tankWidth && float32(math.Abs(float64(y2-y1))) < tankHeight {
				return true
			}
		}
	}
	return false
}
