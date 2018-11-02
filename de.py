import random


def tour(pos):
    t = [(x, i) for (i, x) in enumerate(pos)]
    t.sort()
    return t


def cost(matrix, pos):
    t = tour(pos)
    return sum(matrix[t[i-1][1]][t[i][1]] for i in range(len(t)))


def shuffle(n):
    x = list(range(n))
    random.shuffle(x)
    return tuple(x)


def de(matrix, F=0.75, CR=0.9, NP=100, G_max=200000):
    n = len(matrix)
    pop  = set()
    while len(pop) < NP:
        pop.add(shuffle(n))
    best = random.sample(pop, 1)[0]
    best_cost = cost(matrix, best)
    G = 0

    while len(pop) >= 4 and G < G_max:
        G += 1
        u = list(pop)
        if G % 100 == 0:
            print(G, best_cost, len(pop))
        for x in u:
            if len(pop) < 4: break
            if x not in pop: continue
            a, b, c = random.sample(pop - {x}, 3)
            y = tuple(a[i] + F * (b[i] - c[i]) if random.random() < CR else x[i]
                      for i in range(n))
            # swap if needed
            y_cost = cost(matrix, y)
            x_cost = cost(matrix, x)
            if x_cost >= y_cost:
                pop.discard(x)
                pop.add(y)
                if y_cost < best_cost or random.random() < 0.005:
                    if y_cost < best_cost:
                        best = y
                        best_cost = y_cost
                    if y_cost < cost(matrix, a): pop.discard(a)
                    if y_cost < cost(matrix, b): pop.discard(b)
                    if y_cost < cost(matrix, c): pop.discard(c)

    return [city for _,city in tour(best)], best_cost
