from itertools import repeat
import numpy as np
import matplotlib.patches as mpatches
import matplotlib.pyplot as plt


ALPHA = 0.9999
EPS   = 0.0001


def plot():
    T = np.array(range(1000000))
    G = 25 * (np.ones(1000000) * ALPHA) ** T
    S = 25 / (1 + T * EPS)
    S2 = 25 / (1 + np.exp(T * EPS))

    plt.grid(True)
    plt.semilogx(T, G, color='blue')
    plt.semilogx(T, S, color='red')
    plt.semilogx(T, S2, color='green')
    plt.ylabel('Temperature')
    plt.xlabel('t')
    plt.legend(handles=[
        mpatches.Patch(color='red', label=r'$ T_k = \frac{T_0}{1+k\epsilon} $'),
        mpatches.Patch(color='blue', label=r'$ T_k = \alpha T_{k-1} $'),
        mpatches.Patch(color='green', label=r'$ T_k = \frac{T_0}{1+\exp(k\epsilon)} $'),
    ])
    plt.show()


if __name__ == '__main__':
    plot()
