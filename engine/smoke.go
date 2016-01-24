package engine

type Smoke struct {
	show     bool
	position []*Position
}

func (s *Server) smokeAdd(x, y float32) {
	pos := &Position{x: x, y: y}
	s.smoke.position = append(s.smoke.position, pos)
	s.smoke.show = true
}

func (s *Server) smokeRead() {
	s.smoke.show = false
	s.smoke.position = []*Position{}
}
