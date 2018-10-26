package main

import "flag"
import "time"
import "os"
import "encoding/json"
import "math"
import "math/rand"

type Matrix = [][]int

type Path struct {
	Tour []int `json:"tour"`
	Cost int   `json:"cost"`
}

func neighbour(s []int) []int {
	x := make([]int, len(s))
	copy(x, s)
	i := rand.Intn(len(s) - 1)
	j := (i + 1) + rand.Intn(len(s)-(i+1))
	m := (j - i) / 2
	for k := 0; k < m; k++ {
		x[i+k], x[j-k] = x[j-k], x[i+k]
	}
	return x
}

func cost(s []int, matrix Matrix) int {
	d := 0
	for i := 1; i < len(s); i++ {
		d += matrix[s[i-1]][s[i]]
	}
	d += matrix[s[len(s)-1]][s[0]]
	return d
}

func p(e, next_e, temp float64) float64 {
	return math.Exp((e - next_e) / temp)
}

func initial(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

func anneal(matrix Matrix, alpha float64) ([]int, int) {
	// TODO: we know the memory requirements beforehand (we only need 3
	// arrays, s, s', and s*). So we can implement a zero-alloc version
	// to completely subvert GC.
	n := len(matrix)
	s := initial(n)
	rand.Shuffle(n, func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})

	T_min := 0.00001
	T := float64(n * n)
	e := float64(cost(s, matrix))
	best_s := s
	best_e := e

	for T > T_min {
		for i := 0; i < 100; i++ {
			next_s := neighbour(s)
			next_e := float64(cost(next_s, matrix))
			// if next_e < best_e then necessarily we have r < p(...)
			if next_e < best_e {
				best_s = next_s
				best_e = next_e
				s = next_s
				e = next_e
				continue
			}
			if rand.Float64() < p(e, next_e, T) {
				s = next_s
				e = next_e
			}
		}
		T *= alpha
	}

	return best_s, int(best_e)
}

func main() {
	alphaPtr := flag.Float64("alpha", 0.99670, "T *= alpha")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	// actually do annealing
	tour, cost := anneal(matrix, *alphaPtr)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
