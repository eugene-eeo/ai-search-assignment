import numpy as np
import matplotlib.pyplot as plt
from itertools import product
from json import loads


def parse(f):
    costs = []
    times = []
    for line in f:
        row = loads(line)
        costs.append(row['cost'])
        times.append(row['time'])
    return costs, times


K = [1.50, 1.59, 1.60, 1.70]
M = [2,    1.50, 1.00, 0.89]


def plot(size):
    d = {}
    for k, m in product(K, M):
        d[k,m] = None
        try:
            with open(f'results-{size}-{k}-{m}.txt', mode='r') as fp:
                costs, times = parse(fp)
                d[k,m] = costs, times
        except FileNotFoundError:
            pass

    fig, axs = plt.subplots(4, 4, sharey=True, tight_layout=True, squeeze=True)
    for i, k in enumerate(K):
        for j, m in enumerate(M):
            if d[k,m] is None:
                continue
            costs = d[k,m][0]
            to_plot = np.array(costs) / min(costs)
            ax = axs[i,j]
            ax.hist(to_plot, weights=np.zeros_like(to_plot) + 1. / to_plot.size)
            ax.set_title(f'k={k},m={m},min={min(costs)}', fontsize=6)
            ax.grid(True)
            for tick in ax.xaxis.get_major_ticks(): tick.label.set_fontsize(6)
            for tick in ax.yaxis.get_major_ticks(): tick.label.set_fontsize(6)

    plt.savefig(f'plot-{size}.pdf', dpi=300)


if __name__ == '__main__':
    import sys
    for arg in sys.argv[1:]:
        plot(arg)
