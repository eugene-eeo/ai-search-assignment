import math
import random


def no_opt(matrix, route):
    return route


def two_opt(matrix, route):
    best = route
    improved = True
    while improved:
        improved = False
        for i in range(1, len(route) - 2):
            for j in range(i+1, len(route)):
                if j - i == 1: continue
                new_route = route[:]
                new_route[i:j] = route[j-1:i-1:-1]
                if cost(matrix, new_route) < cost(matrix, best):
                    best = new_route
                    improved = True
        route = best
    return best


def cost(matrix, tour):
    d = 0
    for i in range(len(tour)):
        d += matrix[tour[i-1]][tour[i]]
    return d


def ant(src, matrix, pheromone, beta, p, t0, p0):
    tour = [src]
    weights = {u: 0 for u in range(len(matrix)) if u != src}
    for _ in range(1, len(matrix)):
        total = 0
        for city in weights:
            weight = pheromone[src][city] / (matrix[src][city] ** beta)
            total += weight
            weights[city] = weight
        # greedily pick next city
        if random.random() < p0:
            dst = max(weights, key=weights.__getitem__)
        else:
            r = random.random() * total
            for dst, weight in weights.items():
                r -= weight
                if r < 0:
                    break
        tour.append(dst)
        del weights[dst]
        # evaporate used trail
        pheromone[src][dst] *= (1 - p)
        pheromone[src][dst] += p * t0
        src = dst
    return tour


def greedy(matrix):
    tour = [0]
    src = 0
    for _ in range(1, len(matrix)):
        _, dst = min((cost, i) for i, cost in enumerate(matrix[src]) if i not in tour)
        src = dst
        tour.append(dst)
    return tour


def aco(matrix, G, beta=2, p=0.1, p0=0.9, opt=two_opt):
    n = len(matrix)
    best = list(range(n))
    random.shuffle(best)
    best_cost = cost(matrix, best)
    t0 = 1 / (n * cost(matrix, greedy(matrix)))

    pheromone = [[0] * n for _ in range(n)]
    for i, row in enumerate(matrix):
        for j, distance in enumerate(row):
            pheromone[i][j] = t0 if i != j else None

    for _ in range(G):
        # simulate ants
        for i in range(n):
            tour = ant(i,
                       matrix,
                       pheromone,
                       beta=beta,
                       p=p,
                       t0=t0,
                       p0=p0)
            u = cost(matrix, tour)
            if u < best_cost:
                best = tour
                best_cost = u
        # local search phase
        best = opt(matrix, best)
        best_cost = cost(matrix, best)
        # best ant updates pheromone
        for i, x in enumerate(best):
            y = best[i-1]
            pheromone[x][y] = (1 - p) * pheromone[x][y] + p / best_cost
            pheromone[y][x] = (1 - p) * pheromone[y][x] + p / best_cost

    return best, best_cost
