package main

import "flag"
import "time"
import "os"
import "encoding/json"
import "math"
import "math/rand"

type Matrix = [][]int

type Path struct {
	Tour []int   `json:"tour"`
	Cost int     `json:"cost"`
	Time float64 `json:"time"`
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

func anneal(matrix Matrix, m, k float64) ([]int, int) {
	n := len(matrix)
	s := initial(n)
	rand.Shuffle(n, func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})

	// sometimes 1/(1+eps)^n ~= 0 (floating point)
	alpha := math.Min(1-1/math.Pow(k, float64(n)), 0.9999)
	T_min := 0.00001
	T := m * float64(n)
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
	kPtr := flag.Float64("k", 1.59, "α = 1/k^n")
	mPtr := flag.Float64("m", 1.00, "T = m*n")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	r := json.NewDecoder(os.Stdin)
	err := r.Decode(&matrix)
	if err != nil {
		panic(err)
	}

	// actually do annealing
	t0 := time.Now()
	tour, cost := anneal(matrix, *mPtr, *kPtr)
	t1 := time.Now()

	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
		Time: t1.Sub(t0).Seconds(),
	})
}
