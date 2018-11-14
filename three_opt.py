#!/usr/bin/env python
import json
import sys


def cost(T, M):
    d = 0
    for x, y in T:
        d += M[x][y]
    return d


class Found(Exception):
    def __init__(self, tour):
        self.tour = tour


# Tries to find a series of swaps on T that will improve the tour
# uses a simplified version of the Lin-Kernighan heuristic, while
# only backtracking for the first few depths. Adapted from Lin-Kernighan
# Heuristic Adaptations for the GTSP [https://arxiv.org/pdf/1003.5330.pdf].
def improve(T, M, P, depth, R, maxdepth=5):
    print(depth, R)
    # P = b - ... - x - y - ... - e
    b = P[0][0]
    e = P[-1][-1]

    def gain(edge):
        x, y = edge
        return M[x][y] - M[e][x]

    if depth < maxdepth:
        for i, (x, y) in enumerate(P):
            if x not in R and gain((x, y)) > 0:
                # Q = b - ... - x - e - ... - y - b
                Q = P[:i] + [(x,e)] + [(v, u) for u, v in reversed(P[i+1:])] + [(y,b)]
                if cost(T, M) > cost(Q, M):
                    raise Found(Q)
                improve(T, M, Q[:-1], depth+1, R | {x})
    else:
        x, y = max(P, key=gain)
        if gain((x, y)) > 0:
            i = P.index((x, y))
            # Q = b - ... - x - e - ... - y - b
            Q = P[:i] + [(x, e)] + [(v, u) for u, v in reversed(P[i+1:])] + [(y, b)]
            if cost(T, M) > cost(Q, M):
                raise Found(Q)
            # stop backtracking + proceed greedily
            improve(T, M, Q[:-1], depth+1, R | {x})
            return


def tour2path(T):
    P = []
    for i in range(len(T)):
        x = T[i]
        y = T[(i + 1) % len(T)]
        P.append((x, y))
    return P


def path2tour(P):
    T = []
    for x, y in P:
        T.append(x)
    return T


def main():
    M = json.load(open(sys.argv[1]))
    T = json.load(open(sys.argv[2]))["tour"]
    T = tour2path(T)
    P = T[:-1]
    try:
        improve(T, M, P, 0, set(), maxdepth=10)
    except Found as exc:
        T = exc.tour
    print(json.dumps({
        "tour": path2tour(T),
        "cost": cost(T, M),
    }, indent=None))


if __name__ == '__main__':
    main()
