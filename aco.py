import random


def cost(matrix, tour):
    c = 0
    for i in range(1, len(tour)):
        # includes tour[-1], tour[0]
        c += matrix[tour[i-1]][tour[i]]
    c += matrix[tour[0]][tour[-1]]
    return c


def zero(m):
    n = len(m)
    for i in range(n):
        row = m[i]
        for j in range(n):
            row[j] = 0


def aco(matrix, G=100, num_ants=100):
    alpha = 0.35   # evaporation rate
    explore = 15   # weight for explore
    exploit = 15   # weight for exploit

    n = len(matrix)
    P = [[exploit] * n for _ in range(n)]
    S = [[0] * n for _ in range(n)]
    best = list(range(n))
    best_cost = cost(matrix, best)

    while G > 0:
        G -= 1

        for _ in range(num_ants):
            i = random.randint(0, n - 1)
            tour = [i]
            unvisited = list(range(n))
            unvisited.remove(i)
            while unvisited:
                prob = {
                    j: P[i][j] + explore * 1/matrix[i][j]
                    for j in unvisited
                }
                s = sum(prob.values())
                r = random.random()
                for j, p in prob.items():
                    r -= p / s
                    if r < 0:
                        tour.append(j)
                        unvisited.remove(j)
                        i = j
                        break
            t = cost(matrix, tour)
            for i in range(n):
                S[tour[i-1]][tour[i]] += exploit / t
                S[tour[i]][tour[i-1]] += exploit / t
            if t < best_cost:
                best = tour
                best_cost = t

        for i in range(n):
            for j in range(n):
                P[i][j] *= 1 - alpha
                P[i][j] += S[i][j]
                S[i][j] = 0

    return best, best_cost
