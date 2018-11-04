import heapq
import sys
import os


def cost(matrix, tour):
    if len(tour) < len(matrix):
        return sum(matrix[tour[i-1]][tour[i]] for i in range(1, len(tour)))
    return sum(matrix[tour[i-1]][tour[i]] for i in range(len(tour)))


g = cost


def greedy_search(M, z):
    m = len(z)
    n = len(M)
    if n == m:
        return 0
    V = set(z)
    c = 0
    u = z[-1]
    while len(V) < n:
        v = min((i for i in range(n) if i not in V),
                key=lambda i: M[u][i])
        c += M[u][v]
        u = v
        V.add(v)
    c += M[v][z[0]]
    return c


def a_star(M, h=greedy_search):
    n = len(M)
    t = (0,)
    pq = [(g(M, t) + h(M, t), t)]
    heapq.heapify(pq)
    seen = set()

    while True:
        sys.stderr.write("%s %s\n" % (len(pq), pq[0][0]))
        cost, tour = heapq.heappop(pq)
        if len(tour) == n:
            return cost, tour

        if tour in seen:
            continue

        seen.add(tour)
        visited = set(tour)
        for i in range(n):
            if i not in visited:
                t = tour + (i,)
                if t not in seen:
                    f = g(M, t) + h(M, t)
                    heapq.heappush(pq, (f, t))


if __name__ == '__main__':
    import json
    from read_file import read

    for f in os.listdir('city_data/'):
        sys.stderr.write(f + "\n")
        c, t = a_star(read(os.path.join('city_data/', f)))
        print(json.dumps({
            "name": f,
            "tour": t,
            "cost": c,
            }))
