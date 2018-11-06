package main

import "fmt"
import "flag"
import "os"
import "encoding/json"
import "math"
import "math/rand"

func cost(matrix [][]int, tour []int) int {
	d := 0
	for i := 1; i < len(tour); i++ {
		d += matrix[tour[i-1]][tour[i]]
	}
	d += matrix[tour[len(tour)-1]][tour[0]]
	return d
}

func reverse(x []int, i, j int) {
	for left, right := i, j; left < right; left, right = left+1, right-1 {
		x[left], x[right] = x[right], x[left]
	}
}

func two_opt(matrix [][]int, tour []int) {
	tour_cost := cost(matrix, tour)
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
				c := cost(matrix, route)
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
}

func max_weight(weights map[int]float64) int {
	best := 0
	best_weight := 0.0
	for city, weight := range weights {
		if weight > best_weight {
			best = city
			best_weight = weight
		}
	}
	return best
}

func choose_weighted(weights map[int]float64, total float64) int {
	for {
		r := rand.Float64() * total
		for dst, weight := range weights {
			r -= weight
			if r < 0 {
				return dst
			}
		}
	}
}

func ant(
	src int,
	tour []int, weights map[int]float64, // can be shared with other ants
	matrix [][]int, pheromone [][]float64, // problem specific components
	beta, p_greedy, t0, rho float64, // parameters
) []int {
	for i := 0; i < len(matrix); i++ {
		if i != src {
			weights[i] = 0
		}
	}
	tour[0] = src
	delete(weights, src)
	for i := 1; i < len(matrix); i++ {
		total := 0.0
		for city, _ := range weights {
			weight := pheromone[src][city] / math.Pow(float64(matrix[src][city]), 2)
			total += weight
			weights[city] = weight
		}
		dst := src
		if rand.Float64() < p_greedy {
			dst = max_weight(weights)
		} else {
			dst = choose_weighted(weights, total)
		}
		tour[i] = dst
		delete(weights, dst)
		pheromone[src][dst] = (1-rho)*pheromone[src][dst] + rho*t0
		pheromone[dst][src] = (1-rho)*pheromone[dst][src] + rho*t0
		src = dst
	}
	return tour
}

func nearest_neighbour(matrix [][]int) []int {
	taken := map[int]bool{0: true}
	tour := []int{0}
	for i := 1; i < len(matrix); i++ {
		src := len(tour) - 1
		dst := src
		best := -1
		for j, distance := range matrix[src] {
			if !taken[j] && (best == -1 || distance < best) {
				best = distance
				dst = j
			}
		}
		tour = append(tour, dst)
		taken[dst] = true
	}
	return tour
}

func aco(matrix [][]int, G int, beta float64, rho float64, p_greedy float64) ([]int, int) {
	n := len(matrix)
	m := n / 2
	best := nearest_neighbour(matrix)
	best_cost := cost(matrix, best)
	t0 := 1 / float64(n*best_cost)

	// Create pheromone matrix with initial value t0
	pheromone := make([][]float64, n)
	for i := 0; i < n; i++ {
		pheromone[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if j == i {
				pheromone[i][j] = 0
			} else {
				pheromone[i][j] = t0
			}
		}
	}

	it_best := make([]int, n)
	it_best_cost := best_cost
	copy(it_best, best)

	tour := make([]int, n)
	weights := make(map[int]float64, n)
	for G > 0 {
		if G%10 == 0 {
			fmt.Fprintln(os.Stderr, G, best_cost, it_best_cost)
		}
		G--
		for src := 0; src < m; src++ {
			ant(src,
				tour, weights,
				matrix, pheromone,
				beta, p_greedy, t0, rho)
			u := cost(matrix, tour)
			if u < best_cost {
				copy(best, tour)
				best_cost = u
			}
			if src == 0 || u < it_best_cost {
				copy(it_best, tour)
				it_best_cost = u
			}
		}
		two_opt(matrix, it_best)
		it_best_cost = cost(matrix, it_best)
		if it_best_cost < best_cost {
			copy(best, it_best)
			best_cost = it_best_cost
		}
		gb := it_best
		bc := it_best_cost
		// use global best instead
		if G%20 == 0 {
			gb = best
			bc = best_cost
		}
		for i := 0; i < n; i++ {
			x := gb[i]
			y := gb[(i+1)%n]
			pheromone[x][y] = (1-rho)*pheromone[x][y] + rho/float64(bc)
			pheromone[y][x] = (1-rho)*pheromone[y][x] + rho/float64(bc)
		}
	}

	return best, best_cost
}

type Path struct {
	Tour []int `json:"tour"`
	Cost int   `json:"cost"`
}

func main() {
	betaPtr := flag.Float64("beta", 2.0, "exploration")
	p_greedyPtr := flag.Float64("pg", 0.9, "greedy probability")
	rhoPtr := flag.Float64("rho", 0.9, "evaporation rate")
	GPtr := flag.Int("G", 2000, "iterations")
	flag.Parse()

	//rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	tour, cost := aco(matrix, *GPtr, *betaPtr, *rhoPtr, *p_greedyPtr)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
