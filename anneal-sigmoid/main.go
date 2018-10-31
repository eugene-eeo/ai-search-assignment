package main

import "flag"
import "fmt"
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

func neighbour(x []int, s []int) {
	copy(x, s)
	// Observations:
	//  - cost((1,2,3)) = cost((2,3,1)) = cost((3,1,2)), so first city can be fixed.
	//  - cost((1,2,3)) = cost((3,2,1)), so reverse order is the same cost.
	n := len(s)
	i := 1
	j := n - 1
	for i == 1 && j == n-1 {
		i = 1 + rand.Intn(n-2)
		j = (i + 1) + rand.Intn(n-(i+1))
	}
	m := (j - i) / 2
	for k := 0; k < m; k++ {
		x[i+k], x[j-k] = x[j-k], x[i+k]
	}
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

func ccopy(x []int) []int {
	s := make([]int, len(x))
	copy(s, x)
	return s
}

func anneal(matrix Matrix, epsilon float64) ([]int, int) {
	n := len(matrix)
	s := initial(n)
	rand.Shuffle(n-1, func(i, j int) {
		s[i+1], s[j+1] = s[j+1], s[i+1]
	})

	T_min := 1 - math.Pow(epsilon, 2)
	T_0 := float64(n)
	k := epsilon
	T := T_0
	g := 0
	e := float64(cost(s, matrix))
	next_s := ccopy(s)
	best_s := ccopy(s)
	best_e := e

	for T > T_min {
		g++
		for i := 0; i < 100; i++ {
			neighbour(next_s, s)
			next_e := float64(cost(next_s, matrix))
			// if next_e < best_e then necessarily we have r < p(...)
			if next_e < best_e {
				copy(best_s, next_s)
				best_e = next_e
				copy(s, next_s)
				e = next_e
				continue
			}
			if rand.Float64() < p(e, next_e, T) {
				copy(s, next_s)
				e = next_e
			}
		}
		if g%n == 0 {
			fmt.Fprintln(os.Stderr, T, best_e, e)
		}
		T = T_0 / (1 + k)
		k += epsilon
	}

	return best_s, int(best_e)
}

func main() {
	epsPtr := flag.Float64("epsilon", 0.0001, "T_k = T_0/(1 + e^(-k*epsilon))")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	// actually do annealing
	tour, cost := anneal(matrix, *epsPtr)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
