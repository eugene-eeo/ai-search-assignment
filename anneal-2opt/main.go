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

func reverse(x []int, i, j int) {
	for i < j {
		x[i], x[j] = x[j], x[i]
		i++
		j--
	}
}

func two_opt(tour []int, matrix [][]int) int {
	tour_cost := cost(tour, matrix)
	improved := true
	for improved {
		improved = false
		for i := 1; i < len(tour)-2; i++ {
			for j := i + 1; j < len(tour); j++ {
				if j-i == 1 {
					continue
				}
				reverse(tour, i, j)
				c := cost(tour, matrix)
				if c < tour_cost {
					tour_cost = c
					improved = true
				} else {
					reverse(tour, i, j)
				}
			}
		}
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

func anneal(matrix Matrix, alpha float64) ([]int, int) {
	s := initial(len(matrix))
	rand.Shuffle(len(matrix)-1, func(i, j int) {
		s[i+1], s[j+1] = s[j+1], s[i+1]
	})

	// energies and states
	g := 0
	e := float64(cost(s, matrix))
	next_s := ccopy(s)
	best_s := ccopy(s)
	best_e := float64(cost(s, matrix))

	T := float64(len(matrix))
	//T := 0.0
	//for i := 0; i < len(matrix); i++ {
	//	neighbour(next_s, s)
	//	c := float64(cost(next_s, matrix))
	//	if c > T {
	//		T = c
	//	}
	//}

	// temperature
	k := 0.0
	eps := 1 - alpha
	T_min := 0.0001
	T0 := T

	for T > T_min {
		g++
		fmt.Fprintln(os.Stderr, T, best_e)
		neighbour(next_s, s)
		next_e := float64(two_opt(next_s, matrix))
		if next_e < best_e {
			copy(best_s, next_s)
			best_e = next_e
		}
		if next_e < e || rand.Float64() < p(e, next_e, T) {
			copy(s, next_s)
			e = next_e
		}
		// Geometric schedule
		if T < 5 {
			T *= alpha
		} else {
			k += eps
			T = T0 / (1 + k)
		}
	}

	best_e = float64(two_opt(best_s, matrix))
	return best_s, int(best_e)
}

func main() {
	alphaPtr := flag.Float64("alpha", 0.99670, "T *= alpha")
	//fPtr := flag.Int("f", 10, "debug frequency")
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
