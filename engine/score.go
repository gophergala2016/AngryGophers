package engine

type Scores struct {
	change bool
	client map[int]int
}

func (s *Server) scoreAdd(clientId int) {
	s.score.change = true
	s.score.client[clientId]++
}

func (s *Server) scoreRead() {
	s.score.change = false
}
