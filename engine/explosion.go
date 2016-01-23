package engine

type Explosion struct {
	show     bool
	position []*Position
}

type Position struct {
	x float32
	y float32
}

func (s *Server) explosionAdd(x, y float32) {
	pos := &Position{x: x, y: y}
	s.explosion.position = append(s.explosion.position, pos)
	s.explosion.show = true
}

func (s *Server) explosionRead() {
	s.explosion.show = false
	s.explosion.position = []*Position{}
}
