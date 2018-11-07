import sys
import json
import statistics


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

    print("")


if __name__ == '__main__':
    main()
