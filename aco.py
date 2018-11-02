import math
import random


def cost(matrix, tour):
    d = 0
    for i in range(len(tour)):
        d += matrix[tour[i-1]][tour[i]]
    return d


def ant(src, matrix, pheromone, alpha, beta, epsilon, t0):
    tour = [src]
    weights = {u: 0 for u in range(len(matrix)) if u != src}
    for _ in range(1, len(matrix)):
        greedy = random.random() < 0.3
        a = 1 if greedy else alpha
        total = 0
        for city in weights:
            weight = (pheromone[src][city] ** a) / (matrix[src][city] ** beta)
            total += weight
            weights[city] = weight
        # greedily pick next city
        if greedy:
            dst = max(weights, key=weights.__getitem__)
        else:
            r = random.random() * total
            for dst, weight in weights.items():
                r -= weight
                if r < 0:
                    break
        tour.append(dst)
        del weights[dst]
        pheromone[src][dst] *= (1 - epsilon)
        pheromone[src][dst] += epsilon * t0
        src = dst
    return tour


def aco(matrix, G, alpha=1, beta=2, epsilon=0.1, p=0.2):
    n = len(matrix)
    best = list(range(n))
    random.shuffle(best)
    best_cost = cost(matrix, best)
    t0 = 50

    pheromone = [[0] * n for _ in range(n)]
    for i, row in enumerate(matrix):
        for j, distance in enumerate(row):
            pheromone[i][j] = t0

    for _ in range(G):
        # simulate ants
        for i in range(n):
            tour = ant(i,
                       matrix,
                       pheromone,
                       alpha=alpha,
                       beta=beta,
                       epsilon=epsilon,
                       t0=t0)
            u = cost(matrix, tour)
            if u < best_cost:
                best = tour
                best_cost = u
        # update pheromone
        for i, x in enumerate(best):
            y = best[i-1]
            pheromone[x][y] = (1 - p) * pheromone[x][y] + p / best_cost
            pheromone[y][x] = (1 - p) * pheromone[y][x] + p / best_cost
        print(best_cost)

    return best, best_cost
