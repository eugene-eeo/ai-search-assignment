import re
import glob
import numpy as np
import matplotlib.pyplot as plt
from json import loads


def parse(f):
    costs = []
    for line in f:
        row = loads(line)
        costs.append(row['cost'])
    return costs


def plot(size):
    d = {}
    K = []
    for filename in glob.glob(f'results-{size}-*.txt'):
        k = filename[len(f'results-{size}-'):-len('.txt')]
        K.append(k)
        with open(filename, mode='r') as fp:
            costs = parse(fp)
            d[k] = costs

    min_cost = min(min(costs) for costs in d.values())
    fig, axs = plt.subplots(2, 4, sharey=True, tight_layout=True, squeeze=True)
    for i in range(4):
        for j in range(2):
            k = K[i + 4*j]
            to_plot = np.array(d[k]) / min_cost

            ax = axs[j,i]
            ax.hist(to_plot, weights=np.zeros_like(to_plot) + 1. / to_plot.size)
            ax.grid(True)

            # adjust font
            ax.set_title(f'k={k}', fontsize=6)
            for tick in ax.xaxis.get_major_ticks(): tick.label.set_fontsize(6)
            for tick in ax.yaxis.get_major_ticks(): tick.label.set_fontsize(6)

    plt.savefig(f'plot-{size}.pdf', dpi=300)


if __name__ == '__main__':
    import sys
    for arg in sys.argv[1:]:
        plot(arg)
