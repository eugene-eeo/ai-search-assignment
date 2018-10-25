package main

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

func anneal(matrix Matrix) ([]int, int) {
	n := len(matrix)
	s := initial(n)
	rand.Shuffle(n, func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})

	T := 1.0 + 0.5*math.Log10(float64(n))
	T_min := 0.00001
	alpha := 0.99
	e := float64(cost(s, matrix))

	for T > T_min {
		for i := 0; i < 100; i++ {
			next_s := neighbour(s)
			next_e := float64(cost(next_s, matrix))
			if rand.Float64() < p(e, next_e, T) {
				s = next_s
				e = next_e
			}
		}
		T *= alpha
	}

	return s, int(e)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	r := json.NewDecoder(os.Stdin)
	err := r.Decode(&matrix)
	if err != nil {
		panic(err)
	}
	tour, cost := anneal(matrix)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
