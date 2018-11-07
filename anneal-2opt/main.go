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

func two_opt(matrix [][]int, tour []int) int {
	tour_cost := cost(tour, matrix)
	route := make([]int, len(tour))
	improved := true
	for improved {
		improved = false
		for i := 1; i < len(tour)-2; i++ {
			for j := i + 1; j < len(tour); j++ {
				if j-i == 1 {
					continue
				}
				copy(route, tour)
				reverse(route, i, j)
				c := cost(route, matrix)
				if c < tour_cost {
					tour_cost = c
					copy(tour, route)
					improved = true
				}
			}
		}
		// Iterate on best found
		copy(route, tour)
	}
	return tour_cost
}

func neighbour(x []int, s []int) {
	copy(x, s)
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
	s := initial(len(matrix))
	rand.Shuffle(len(matrix)-1, func(i, j int) {
		s[i+1], s[j+1] = s[j+1], s[i+1]
	})

	// energies and states
	e := float64(cost(s, matrix))
	next_s := ccopy(s)
	best_s := ccopy(s)
	best_e := float64(cost(s, matrix))

	T := 0.0
	for i := 0; i < len(matrix); i++ {
		neighbour(next_s, s)
		c := float64(cost(next_s, matrix))
		if c > T {
			T = c
		}
	}

	// temperature
	g := 0
	k := 0.0
	eps := 1 - alpha
	T_min := math.Pow(1-eps, 2)
	T0 := T

	for T > T_min {
		if g%debugFreq == 0 {
			fmt.Fprintln(os.Stderr, g, T, best_e)
		}
		if rand.Float64() < 1/math.Pow(T, 2.0) {
			copy(s, best_s)
			e = best_e
		}
		g++
		for i := 0; i < 10; i++ {
			neighbour(next_s, s)
			next_e := float64(two_opt(matrix, next_s))
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
		if T < 5 {
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
