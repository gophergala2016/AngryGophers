package engine

import (
	"log"
	"math/rand"
	"time"
)

type mapa struct {
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

func (s *mapa) drawMap() ([][]int, []int) {
	return s.ground, s.speedGround
}

func (s *mapa) GetTrees() [][]int {
	return s.trees
}

func (s *mapa) GetRocks() [][]int {
	return s.rocks
}

func getMap(ground [][]int, speedGround []int, mapSizeX, mapSizeY float32) *mapa {
	s := &mapa{}
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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	cntTree := r.Intn(20)
	for i := 0; i < cntTree; i++ {
		x := (r.Intn(int(mapSizeX)))
		y := (r.Intn(int(mapSizeY)))
		s.trees = append(s.trees, []int{x, y})
	}

	cntRock := r.Intn(20)
	for i := 0; i < cntRock; i++ {
		x := (r.Intn(int(mapSizeX)))
		y := (r.Intn(int(mapSizeY)))
		s.rocks = append(s.rocks, []int{x, y})
	}

	// // woda
	// for x := 280; x < 380; x++ {
	// 	for y := 580; y < 800; y++ {
	// 		res[x][y] = 3
	// 	}
	// }
	// // lod
	// for x := 280; x < 380; x++ {
	// 	for y := 280; y < 480; y++ {
	// 		res[x][y] = 30
	// 	}
	// }
	// // lod
	// for x := 380; x < 480; x++ {
	// 	for y := 280; y < 380; y++ {
	// 		res[x][y] = 30
	// 	}
	// }
	// // drzewo
	// for x := 167; x < 396; x++ {
	// 	for y := 67; y < 200; y++ {
	// 		res[x][y] = 0
	// 	}
	// }
	// // drzewo
	// for x := 267; x < 396; x++ {
	// 	for y := 0; y < 67; y++ {
	// 		res[x][y] = 0
	// 	}
	// }
	s.speedPoint = res
	return s
}
