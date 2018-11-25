import sys
import json


def main():
    min_cost = float('+inf')
    min_tour = []

    for i, line in enumerate(sys.stdin):
        d = json.loads(line)
        cost = d['cost']
        tour = d['tour']
        if cost < min_cost:
            min_cost = cost
            min_tour = tour

    n = len(min_tour)
    print("NAME = AISearchfile{:03},".format(n))
    print("TOURSIZE = {},".format(n))
    print("LENGTH = {},".format(min_cost))
    print(",".join(str(k+1) for k in min_tour))


if __name__ == '__main__':
    main()
