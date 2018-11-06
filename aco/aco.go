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
	// Goal: to find the pair of edges (i,i+1) and (j,j+1) s.t.
	// replacing them with (i,j) and (i+1,j+1) produces a tour
	// of smaller cost.
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

func max_weight(infos []*cityInfo) int {
	best := 0
	best_weight := 0.0
	for city, info := range infos {
		if !info.visited && info.weight > best_weight {
			best = city
			best_weight = info.weight
		}
	}
	return best
}

func choose_weighted(infos []*cityInfo, total float64) int {
	if math.IsInf(total, +1) {
		for dst, info := range infos {
			if !info.visited {
				return dst
			}
		}
	}
	for {
		r := rand.Float64() * total
		for dst, info := range infos {
			if !info.visited {
				r -= info.weight
				if r < 0 {
					return dst
				}
			}
		}
	}
}

type cityInfo struct {
	weight  float64
	visited bool
}

func ant(
	tour []int, infos []*cityInfo, // can be shared with other ants
	matrix [][]int, pheromone [][]float64, // problem specific components
	beta, p_greedy, t0, rho float64, // parameters
) []int {
	// initialize tour and infos
	src := rand.Intn(len(matrix))
	tour[0] = src
	for i := 0; i < len(matrix); i++ {
		infos[i].visited = i == src
	}
	for i := 1; i < len(matrix); i++ {
		total := 0.0
		for city, info := range infos {
			if !info.visited {
				info.weight = pheromone[src][city] / math.Pow(float64(matrix[src][city]), beta)
				total += info.weight
			}
		}
		dst := src
		if rand.Float64() < p_greedy {
			dst = max_weight(infos)
		} else {
			dst = choose_weighted(infos, total)
		}
		tour[i] = dst
		infos[dst].visited = true
		pheromone[src][dst] = (1-rho)*pheromone[src][dst] + rho*t0
		pheromone[dst][src] = (1-rho)*pheromone[dst][src] + rho*t0
		src = dst
	}
	two_opt(matrix, tour)
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
	t0 := 1 / (float64(n) * float64(best_cost))

	// Create pheromone matrix with initial value t0
	pheromone := make([][]float64, n)
	for i := 0; i < n; i++ {
		pheromone[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			pheromone[i][j] = t0
		}
	}

	it_best := make([]int, n)
	it_best_cost := best_cost
	copy(it_best, best)

	tour := make([]int, n)
	infos := make([]*cityInfo, n)
	for i := 0; i < n; i++ {
		infos[i] = &cityInfo{}
	}

	for G > 0 {
		fmt.Fprintln(os.Stderr, G, best_cost, it_best_cost)
		G--
		for i := 0; i < m; i++ {
			ant(tour, infos,
				matrix, pheromone,
				beta, p_greedy, t0, rho)
			u := cost(matrix, tour)
			if u < best_cost {
				copy(best, tour)
				best_cost = u
			}
			if i == 0 || u < it_best_cost {
				copy(it_best, tour)
				it_best_cost = u
			}
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
