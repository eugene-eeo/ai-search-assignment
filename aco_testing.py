import json
from itertools import product
from collections import namedtuple
from aco import aco, two_opt, no_opt
from read_file import read


# matrix of configs
G     = [3000]
OPTS  = [two_opt, no_opt]
BETA  = [1, 2, 3]
P     = [0.1, 0.2, 0.3, 0.4, 0.5]
P0    = [0.9, 0.8, 0.7, 0.6, 0.5]
SIZES = ['058']


Params = namedtuple('Params', 'G,beta,p,p0,opt')
PARAMS = [Params(*p) for p in product(G, BETA, P, P0, OPTS)]


if __name__ == '__main__':
    for size in SIZES:
        matrix = read(f'city_data/AISearchFile{size}.txt')
        with open(f'{size}.jsonl', 'a') as fp:
            for params in PARAMS:
                d = dict(params._asdict())
                d['opt'] = 'two_opt' if params.opt is two_opt else 'none'
                print(size, d)
                tour, cost = aco(matrix, **params._asdict())
                fp.write(json.dumps({
                    "params": d,
                    "tour": tour,
                    "cost": cost,
                }))
                fp.write("\n")
                fp.flush()
