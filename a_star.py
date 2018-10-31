import heapq


def cost(matrix, tour):
    return sum(matrix[tour[i-1]][tour[i]] for i in range(len(tour)))


g = cost


def h1(M, z):
    m = len(z)
    n = len(M)
    if n == m:
        return 0
    V = set(z)
    c = 0
    u = z[-1]
    while m < n:
        v = min((i for i in range(n) if i not in V),
                key=lambda i: M[u][i])
        c += M[u][v]
        m += 1
        u = v
        V.add(v)
    c += M[v][z[0]]
    return c


def h2(M, z):
    m = len(z)
    n = len(M)
    if m < n:
        V = set(z)
        return min(cost for i, cost in enumerate(M[z[-1]]) if i not in V)
    return 0


def a_star(M, h=lambda M,t: h1(M,t) + h2(M,t)):
    n = len(M)
    t = (0,)
    pq = [(g(M, t) + h(M, t), t)]
    seen = set()

    while True:
        cost, tour = heapq.heappop(pq)
        if tour in seen:
            continue

        if len(tour) == n:
            return cost, tour

        seen.add(tour)
        visited = set(tour)
        for i in range(n):
            if i not in visited:
                t = tour + (i,)
                f = g(M, t) + h(M, t)
                if t not in seen:
                    heapq.heappush(pq, (f, t))
