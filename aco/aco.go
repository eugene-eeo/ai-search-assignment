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

type searchNode struct {
	tour  []int
	depth int
	used  map[int]bool
}

func union(A map[int]bool, x int) map[int]bool {
	u := make(map[int]bool, len(A)+1)
	for k, v := range A {
		u[k] = v
	}
	u[x] = true
	return u
}

func lk_opt(M [][]int, P []int) []int {
	n := len(M)
	S := []searchNode{}
	S = append(S, searchNode{
		tour:  P,
		depth: 0,
		used:  map[int]bool{},
	})
	var node searchNode
	for len(S) > 0 {
		node, S = S[len(S)-1], S[:len(S)-1]
		// actually improve here
		e := node.tour[n-1]
		if node.depth < MAX_DEPTH {
			for i, x := range node.tour {
				y := node.tour[(i+1)%n]
				if node.used[x] || (M[x][y]-M[e][x] <= 0) {
					continue
				}
				Q := make([]int, n)
				copy(Q, P)
				reverse(Q, i+1, n-1) // reverse from y...e => ...xe...y
				if cost(M, Q) < cost(M, P) {
					return Q
				}
				S = append(S, searchNode{
					tour:  Q,
					depth: node.depth + 1,
					used:  union(node.used, x),
				})
			}
		} else {
			max_g := 0
			max_i := -1
			max_x := -1
			max_y := -1
			for i, x := range node.tour {
				if node.used[x] {
					continue
				}
				y := node.tour[(i+1)%n]
				g := M[x][y] - M[e][x]
				if g > max_g {
					max_g = g
					max_i = i
					max_x = x
					max_y = y
				}
			}
			if max_i != -1 && (M[max_x][max_y]-M[e][max_x] > 0) {
				Q := make([]int, n)
				copy(Q, P)
				reverse(Q, max_i+1, n-1)
				if cost(M, Q) < cost(M, P) {
					return Q
				}
				S = append(S, searchNode{
					tour:  Q,
					depth: node.depth + 1,
					used:  union(node.used, max_x),
				})
			}
		}
	}
	return P
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
	tour []int, infos []*cityInfo, // can be shared with other ants
	matrix [][]int, pheromone [][]float64, // problem specific components
	beta, p_greedy, t0, rho float64, // parameters
) int {
	// initialize tour and infos
	n := len(matrix)
	src := rand.Intn(n)
	tour[0] = src
	for i := 0; i < n; i++ {
		infos[i].visited = i == src
	}
	for i := 1; i < n; i++ {
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
		pheromone[src][dst] *= 1 - rho
		pheromone[src][dst] += rho * t0
		src = dst
	}
	// make sure to update last edge
	pheromone[tour[n-1]][tour[0]] = (1-rho)*pheromone[tour[n-1]][tour[0]] + rho*t0
	// for some reason, updating tour edges used in 2-opt is worse than updating
	// the edges used to produce the tour ('wrong' edges).
	return two_opt(tour, matrix)
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

	tour := make([]int, n)
	infos := make([]*cityInfo, n)
	for i := 0; i < n; i++ {
		infos[i] = &cityInfo{}
	}

	for G > 0 {
		if debug {
			fmt.Fprintln(os.Stderr, G, best_cost, it_best_cost)
		}
		G--
		for i := 0; i < m; i++ {
			tour_cost := ant(tour, infos, matrix, pheromone, beta, p_greedy, t0, rho)
			if tour_cost < best_cost {
				copy(best, tour)
				best_cost = tour_cost
			}
			if i == 0 || tour_cost < it_best_cost {
				copy(it_best, tour)
				it_best_cost = tour_cost
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
			pheromone[x][y] *= 1 - rho
			pheromone[x][y] += rho / float64(bc)
		}
	}
	best = lk_opt(matrix, best)
	best_cost = cost(matrix, best)
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
