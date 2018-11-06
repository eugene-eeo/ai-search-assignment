package main

import "fmt"
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

func reverse(x []int, i, j int) {
	for i < j {
		x[i], x[j] = x[j], x[i]
		i++
		j--
	}
}

func two_opt(matrix [][]int, tour []int) {
	n := len(tour)
	minchange := -1
	for minchange < 0 {
		minchange = 0
		min_i := -1
		min_j := -1
		for i := 0; i < n-2; i++ {
			for j := i + 2; j < n; j++ {
				j_1 := (j + 1) % n
				change := matrix[tour[i]][tour[j]] + matrix[tour[i+1]][tour[j_1]] - matrix[tour[i]][tour[i+1]] - matrix[tour[j]][tour[j_1]]
				if change < minchange {
					minchange = change
					min_i = i
					min_j = j
				}
			}
		}
		if min_i != -1 {
			j_1 := (min_j + 1) % n
			tour[min_i], tour[min_i+1], tour[min_j], tour[j_1] = tour[min_i], tour[min_j], tour[min_i+1], tour[j_1]
		}
	}
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
	reverse(x, i, j)
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

func anneal(matrix Matrix, alpha float64, debugFreq int) ([]int, int) {
	n := len(matrix)
	s := initial(n)
	rand.Shuffle(n-1, func(i, j int) {
		s[i+1], s[j+1] = s[j+1], s[i+1]
	})

	// temperature
	k := 0.0
	eps := 1 - alpha
	T_c := 0.0
	T_min := math.Pow(1-eps, 2)
	N := float64(n)
	T := float64(n * n)

	// energies and states
	e := float64(cost(s, matrix))
	next_s := ccopy(s)
	best_s := ccopy(s)
	best_e := e

	g := 0

	for T > T_min {
		g++
		if g%debugFreq == 0 {
			fmt.Fprintln(os.Stderr, g, T, best_e)
		}
		for i := 0; i < 100; i++ {
			neighbour(next_s, s)
			two_opt(matrix, next_s)
			next_e := float64(cost(next_s, matrix))
			// if next_e < best_e then necessarily we have r() < p(...)
			if next_e < best_e {
				copy(best_s, next_s)
				copy(s, next_s)
				best_e = next_e
				e = next_e
				continue
			}
			if next_e < e || rand.Float64() < p(e, next_e, T) {
				copy(s, next_s)
				e = next_e
			}
		}
		// Geometric schedule
		if T > N || T < 1 {
			T *= alpha
			T_c = T
		} else {
			k += eps
			T = T_c / (1 + k)
		}
	}

	return best_s, int(best_e)
}

func main() {
	alphaPtr := flag.Float64("alpha", 0.99670, "T *= alpha")
	fPtr := flag.Int("f", 10, "debug frequency")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	// actually do annealing
	tour, cost := anneal(matrix, *alphaPtr, *fPtr)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
