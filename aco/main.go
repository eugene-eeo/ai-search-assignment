package main

import "encoding/json"
import "flag"
import "fmt"
import "math"
import "math/rand"
import "os"

type Matrix = [][]int

type Path struct {
	Tour []int `json:"tour"`
	Cost int   `json:"cost"`
}

func pherome_matrix(matrix [][]float64) [][]float64 {
	x := make([][]float64, len(matrix))
	for i, _ := range x {
		x[i] = make([]float64, len(matrix))
		for j, _ := range x[i] {
			if i != j {
				x[i][j] = 1 / matrix[i][j]
			}
		}
	}
	return x
}

func zeroes(n int) [][]float64 {
	x := make([][]float64, n)
	for i := 0; i < n; i++ {
		x[i] = make([]float64, n)
	}
	return x
}

func initial(n int) []int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = i
	}
	return x
}

func cost(matrix [][]float64, tour []int) float64 {
	d := 0.0
	for i := 1; i < len(tour); i++ {
		d += matrix[tour[i-1]][tour[i]]
	}
	d += matrix[tour[len(tour)-1]][tour[0]]
	return d
}

type pair struct {
	visited bool
	weight  float64
}

func aco(matrix [][]float64, G_max int, alpha, explore, exploit float64) ([]int, float64) {
	n := len(matrix)
	P := pherome_matrix(matrix)
	S := zeroes(n)
	Q := 0.2
	num_ants := n
	best_tour := initial(n)
	best_cost := cost(matrix, best_tour)

	tour := make([]int, n)
	weights := make([]pair, n)

	// TODO:
	// parallelise this by exploiting the fact that we can compute ~4 ants tours
	// at once, ant-worker just needs to return tour and cost.
	for G_max > 0 {
		G_max--
		if (G_max % 1000) == 0 {
			fmt.Fprintln(os.Stderr, G_max, best_cost)
		}

		for ant := 0; ant < num_ants; ant++ {
			src := rand.Intn(n)
			tour[0] = src
			for v := 0; v < n; v++ {
				weights[v].visited = (v == src)
			}
			for z := 1; z < n; z++ {
				total := 0.0
				dest := src
				for v, info := range weights {
					if !info.visited {
						// + eps to alleviate 0 weight
						info.weight = math.Pow(P[src][v], exploit)/math.Pow(matrix[src][v], explore) + 0.001
						total += info.weight
						weights[v] = info
					}
				}
				r := rand.Float64()
				x := 0.0
				for v, info := range weights {
					if !info.visited {
						x += info.weight / total
						if r <= x {
							dest = v
							break
						}
					}
				}
				tour[z] = dest
				weights[dest].visited = true
				src = dest
			}
			c := cost(matrix, tour)
			k := Q / c
			for i := 1; i < n; i++ {
				S[tour[i-1]][tour[i]] += k
				S[tour[i]][tour[i-1]] += k
			}
			S[tour[n-1]][tour[0]] += k
			S[tour[0]][tour[n-1]] += k
			if c < best_cost {
				copy(best_tour, tour)
				best_cost = c
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				P[i][j] = (1-alpha)*P[i][j] + S[i][j]
				S[i][j] = 0
			}
		}
	}

	return best_tour, best_cost
}

func main() {
	alphaPtr := flag.Float64("alpha", 0.85, "pherome evaporation rate")
	explorePtr := flag.Float64("explore", 20, "exploration weight")
	exploitPtr := flag.Float64("exploit", 2.5, "exploitation weight")
	flag.Parse()

	matrix := [][]float64{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}
	n := len(matrix)
	tour, cost := aco(
		matrix,
		2*n*n*n,
		*alphaPtr,
		*explorePtr,
		*exploitPtr,
	)
	w := json.NewEncoder(os.Stdout)
	w.Encode(Path{
		Tour: tour,
		Cost: int(cost),
	})
}
