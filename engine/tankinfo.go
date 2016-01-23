package engine

const tankWidth float32 = 37.5
const tankWidthHalf float32 = 18.75
const tankHeight float32 = 35
const tankHeightHalf float32 = 17.5

const defaultTankSpeed float32 = 5

func (s *Server) checkHitTank(c *Client) (bool, int) {
	hit, hitClientId := s.checkHitBullet(c.id, c.PositionX, c.PositionY, c.PositionX+tankWidth, c.PositionY+tankHeight)
	if hit {
		return true, hitClientId
	}
	return false, 0
}

func (s *Server) setSpeedTank(c *Client, refreshTime float32) float32 {
	return c.Speed * refreshTime * s.getSpeedPosition(c.PositionX+tankWidthHalf, c.PositionY+tankHeightHalf)
}

func (s *Server) checkColision(c *Client, newX, newY float32) {
	if s.getCollision(newX+tankWidthHalf, newY+tankHeightHalf) {
		return
	}
	c.PositionX = newX
	c.PositionY = newY
}
