package main

import "flag"
import "fmt"
import "os"
import "encoding/json"
import "math"
import "math/rand"

type Path struct {
	Tour []int `json:"tour"`
	Cost int   `json:"cost"`
}

func reverse(x []int, i, j int) {
	for i < j {
		x[i], x[j] = x[j], x[i]
		i++
		j--
	}
}

func two_opt(V []int, M [][]int) {
	n := len(M)
	improved := true
	for improved {
		improved = false
		for i := 1; i < n-2; i++ {
			for j := i + 2; j < n; j++ {
				change := M[V[i-1]][V[i]] + M[V[j]][V[(j+1)%n]] - M[V[i-1]][V[j]] - M[V[i]][V[(j+1)%n]]
				if change > 0 {
					reverse(V, i, j)
					improved = true
				}
			}
		}
	}
}

func neighbour(x []int, s []int) {
	copy(x, s)
	n := len(s)
	i := 1
	j := 1
	k := 1
	for !(i < j && j < k) {
		i = 1 + rand.Intn(n-1)
		j = 1 + rand.Intn(n-1)
		k = 1 + rand.Intn(n-1)
	}
	reverse(x, i, j)
	reverse(x, j, k)
}

func cost(V []int, M [][]int) int {
	d := 0
	for i := 1; i < len(V); i++ {
		d += M[V[i-1]][V[i]]
	}
	d += M[V[len(V)-1]][V[0]]
	return d
}

func p(e, next_e, temp float64) float64 {
	// assume that next_e > e
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

func anneal(matrix [][]int, alpha float64, debugFreq int, debug bool) ([]int, int) {
	s := initial(len(matrix))
	rand.Shuffle(len(matrix)-1, func(i, j int) {
		s[i+1], s[j+1] = s[j+1], s[i+1]
	})

	// energies and states
	e := float64(cost(s, matrix))
	next_s := ccopy(s)
	best_s := ccopy(s)
	best_e := float64(cost(s, matrix))

	T := float64(len(matrix))
	G := 0

	// temperature
	k := 0.0
	eps := 1 - alpha
	T_min := 0.0001
	T0 := T

	for T > T_min {
		if debug && G%debugFreq == 0 {
			fmt.Fprintln(os.Stderr, T, best_e)
		}
		G++
		for i := 0; i < 5; i++ {
			neighbour(next_s, s)
			two_opt(next_s, matrix)
			next_e := float64(cost(next_s, matrix))
			if next_e < best_e {
				copy(best_s, next_s)
				best_e = next_e
			}
			if next_e < e || rand.Float64() < p(e, next_e, T) {
				copy(s, next_s)
				e = next_e
			}
		}
		// Geometric schedule
		if T < 1 {
			T *= alpha
		} else {
			k += eps
			T = T0 / (1 + k)
		}
	}

	return best_s, int(best_e)
}

func main() {
	alphaPtr := flag.Float64("alpha", 0.99670, "T *= alpha")
	debugPtr := flag.Int("debug", 0, "debug frequency")
	flag.Parse()

	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	// actually do annealing
	tour, cost := anneal(matrix, *alphaPtr, *debugPtr, *debugPtr > 0)
	json.NewEncoder(os.Stdout).Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
