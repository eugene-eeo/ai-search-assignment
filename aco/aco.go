package main

import "fmt"
import "flag"
import "os"
import "encoding/json"
import "math"
import "math/rand"

const MAX_DEPTH int = 5

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

func two_opt(tour []int, matrix [][]int) int {
	n := len(matrix)
	improved := true
	for improved {
		improved = false
		for i := 1; i < n-2; i++ {
			for j := i + 2; j < n; j++ {
				change := matrix[tour[i-1]][tour[i]] + matrix[tour[j]][tour[(j+1)%n]] - matrix[tour[i-1]][tour[j]] - matrix[tour[i]][tour[(j+1)%n]]
				if change > 0 {
					reverse(tour, i, j)
					improved = true
				}
			}
		}
	}
	return cost(matrix, tour)
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
	step int, tour []int, infos []*cityInfo,
	matrix [][]int, pheromone [][]float64, // problem specific components
	beta, p_greedy, t0, rho float64, // parameters
) {
	if step == 0 {
		src := rand.Intn(len(tour))
		tour[0] = src
		for i, info := range infos {
			info.visited = i == src
		}
		return
	}
	src := tour[step-1]
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
	tour[step] = dst
	infos[dst].visited = true
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

func aco(matrix [][]int, G int, beta float64, rho float64, p_greedy float64, debug bool) ([]int, int) {
	n := len(matrix)
	m := 20
	best := nearest_neighbour(matrix)
	best_cost := cost(matrix, best)
	t0 := 1 / float64(n*best_cost)

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

	tours := [20][]int{}
	infos := [20][]*cityInfo{}
	for i := 0; i < m; i++ {
		infos[i] = make([]*cityInfo, n)
		tours[i] = make([]int, n)
		for j := 0; j < n; j++ {
			infos[i][j] = &cityInfo{}
		}
	}

	for G > 0 {
		G--
		for step := 0; step < n; step++ {
			for i := 0; i < m; i++ {
				ant(step, tours[i], infos[i], matrix, pheromone, beta, p_greedy, t0, rho)
			}
			if step > 0 {
				for i := 0; i < m; i++ {
					// don't perform local update yet
					pheromone[tours[i][step-1]][tours[i][step]] *= 1 - rho
					pheromone[tours[i][step-1]][tours[i][step]] += rho * t0
					pheromone[tours[i][step]][tours[i][step-1]] *= 1 - rho
					pheromone[tours[i][step]][tours[i][step-1]] += rho * t0
				}
			}
		}
		for i := 0; i < m; i++ {
			pheromone[tours[i][n-1]][tours[i][0]] *= 1 - rho
			pheromone[tours[i][0]][tours[i][n-1]] *= 1 - rho
			pheromone[tours[i][n-1]][tours[i][0]] += rho * t0
			pheromone[tours[i][0]][tours[i][n-1]] += rho * t0
			two_opt(tours[i], matrix)
		}
		for i, tour := range tours {
			tour_cost := cost(matrix, tour)
			if tour_cost < best_cost {
				copy(best, tour)
				best_cost = tour_cost
			}
			if i == 0 || tour_cost < it_best_cost {
				copy(it_best, tour)
				it_best_cost = tour_cost
			}
		}
		// every 20 rounds use local best
		gb := best
		bc := best_cost
		if G%20 == 0 {
			gb = it_best
			bc = it_best_cost
		}
		for i := 0; i < n; i++ {
			x := gb[i]
			y := gb[(i+1)%n]
			pheromone[x][y] *= 1 - rho
			pheromone[x][y] += rho / float64(bc)
		}
		if debug {
			costs := [20]int{}
			for i, tour := range tours {
				costs[i] = cost(matrix, tour)
			}
			fmt.Fprintln(os.Stderr, G, best_cost, costs[:5])
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
	debugPtr := flag.Bool("debug", false, "debug")
	GPtr := flag.Int("G", 2000, "iterations")
	flag.Parse()

	//rand.Seed(time.Now().UnixNano())
	matrix := [][]int{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}

	tour, cost := aco(matrix, *GPtr, *betaPtr, *rhoPtr, *p_greedyPtr, *debugPtr)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: cost,
	})
}
