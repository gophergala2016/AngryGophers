package engine

import (
	"log"
)

type Mapa struct {
	name        string
	ground      [][]int
	speedPoint  [][]int
	speedGround []int
	trees       [][]int
	rocks       [][]int
}

func (s *Server) getSpeedPosition(x, y float32) float32 {
	// log.Print(x, y)
	return float32(s.mapa.speedPoint[int(x)][int(y)]) / 10
}

func (s *Server) getCollision(x, y float32) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Print(x, y)
		}
	}()
	// log.Print(x, y)
	return s.mapa.speedPoint[int(x)][int(y)] == 0
}

func (s *Mapa) drawMap() ( string) {
	return s.name
}

func (s *Mapa) GetTrees() [][]int {
	return s.trees
}

func (s *Mapa) GetRocks() [][]int {
	return s.rocks
}

func GetMap(ground [][]int, speedGround []int, mapSizeX, mapSizeY float32) *Mapa {
	s := &Mapa{}
	s.name = "1"
	s.ground = ground
	s.speedGround = speedGround

	res := make([][]int, 800)
	for x := 0; x < len(res); x++ {
		res[x] = make([]int, 800)
	}

	for k1, v1 := range s.ground {
		for k2, v2 := range v1 {
			speed := speedGround[v2]
			for j := 0; j < 50; j++ {
				for i := 0; i < 50; i++ {
					res[(k2*50)+j][(k1*50)+i] = speed
				}
			}
		}
	}

	s.trees = TreeList1
	s.rocks = RockList1

	s.speedPoint = res
	return s
}
