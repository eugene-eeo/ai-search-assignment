package main

import "encoding/json"
import "flag"
import "math/rand"
import "os"

type Matrix = [][]int

type Path struct {
	Tour []int `json:"tour"`
	Cost int   `json:"cost"`
}

func make_pheromone_matrix(n int) [][]float64 {
	x := make([][]float64, n)
	for i := 0; i < n; i++ {
		x[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			x[i][j] = 5.0
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

func aco(matrix [][]float64, G_max int, alpha, explore, exploit float64) ([]int, float64) {
	n := len(matrix)
	P := make_pheromone_matrix(n)
	S := zeroes(n)
	num_ants := n
	best_tour := initial(n)
	best_cost := cost(matrix, best_tour)

	// TODO:
	// parallelise this by exploiting the fact that we can compute ~4 ants tours
	// at once, ant-worker just needs to return tour and cost.
	for G_max > 0 {
		G_max--
		tour := make([]int, n)
		weights := make(map[int]float64, n-1)

		for ant := 0; ant < num_ants; ant++ {
			u := rand.Intn(n)
			tour[0] = u
			for v, _ := range matrix[u] {
				if u != v {
					weights[v] = 0
				}
			}
			for z := 1; z < n; z++ {
				total := 0.0
				dest := 0
				p := rand.Float64()
				if p < 0.3 {
					min := 10000000.0
					for city, distance := range matrix[u] {
						_, ok := weights[city]
						if ok && distance < min {
							dest = city
							min = distance
						}
					}
				} else {
					for v, _ := range weights {
						weight := P[u][v] + explore/matrix[u][v]
						weights[v] = weight
						total += weight
					}
					r := rand.Float64()
					for v, weight := range weights {
						r -= weight / total
						if r < 0 {
							dest = v
							break
						}
					}
				}
				tour[z] = dest
				delete(weights, dest)
				u = dest
			}
			c := cost(matrix, tour)
			k := exploit / c
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
	alphaPtr := flag.Float64("alpha", 0.8, "pherome evaporation rate")
	explorePtr := flag.Float64("explore", 20, "exploration weight")
	exploitPtr := flag.Float64("exploit", 15, "exploitation weight")
	flag.Parse()

	matrix := [][]float64{}
	err := json.NewDecoder(os.Stdin).Decode(&matrix)
	if err != nil {
		panic(err)
	}
	n := len(matrix)
	tour, cost := aco(
		matrix,
		n*n*n*n,
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
