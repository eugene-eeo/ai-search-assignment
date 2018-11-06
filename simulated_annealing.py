import math
import random


def neighbour(s):
    x = s[:]
    i = random.randint(0,     len(s) - 2)
    j = random.randint(i + 1, len(s) - 1)
    x[i:j] = reversed(x[i:j])
    return x


def cost(s, matrix):
    d = 0
    for i in range(len(s)):
        d += matrix[s[i-1]][s[i]]
    return d


def P(e, next_e, temp):
    return math.exp((e - next_e) / temp)


def anneal(matrix, K):
    n = len(matrix)
    s = list(range(n))
    random.shuffle(s)
    e = cost(s, matrix)

    B = s
    E = e

    L = [100] * 5
    r = random.random
    k = 0

    while k < K:
        k += 1
        t_max = max(L)
        print(k, t_max, e)
        t = 0
        c = 0
        for _ in range(100):
            next_s = neighbour(s)
            next_e = cost(next_s, matrix)
            if next_e < e:
                s = next_s
                e = next_e
                if next_e < E:
                    B = next_s
                    E = next_e
            else:
                p = r()
                if p < P(e, next_e, t_max):
                    t = (t - (next_e - e)) / math.log(p, math.e)
                    s = next_s
                    e = next_e
                    c += 1
        if c > 0:
            L.remove(t_max)
            L.append(max(t / c, 0.01))
    return B, E
