import sys
import json
import statistics


def main():
    min_cost = float('+inf')
    max_cost = float('-inf')
    min_tour = []
    costs = []

    for line in sys.stdin:
        d = json.loads(line)
        cost = d['cost']
        tour = d['tour']
        costs.append(cost)
        if cost > max_cost:
            max_cost = cost
        if cost < min_cost:
            min_cost = cost
            min_tour = tour

    print("Mean:     ", statistics.mean(costs))
    print("Variance: ", statistics.variance(costs))
    print("Min:      ", min(costs))
    print("Max:      ", max(costs))
    print("Min tour: ", min_tour)


if __name__ == '__main__':
    main()
