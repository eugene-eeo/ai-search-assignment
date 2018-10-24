import math
import random


def neighbour(s):
    x = s[:]
    i = j = 0
    while i == j:
        i = random.randint(0, len(x) - 1)
        j = random.randint(0, len(x) - 1)
    x[i:j] = reversed(x[i:j])
    return x


def cost(s, matrix):
    d = 0
    for i in range(1, len(s)):
        d += matrix[s[i-1]][s[i]]
    return d


def P(e, next_e, temp):
    return math.exp((e - next_e) / temp)


def find_tour(matrix):
    n = len(matrix)
    s = list(range(n))
    random.shuffle(s)

    T = 1.0
    T_min = 0.00001
    alpha = 0.996
    e = cost(s, matrix)

    while T > T_min:
        for _ in range(100):
            next_s = neighbour(s)
            next_e = cost(next_s, matrix)
            if random.random() < P(e, next_e, T):
                s = next_s
                e = next_e
        T *= alpha
    return s, e
