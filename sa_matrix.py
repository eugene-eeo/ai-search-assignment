import json
import subprocess

K = [
    1.40,
    1.50,
    1.59,
    1.60,
    1.70,
]

M = [
    2,
    1.75,
    1.50,
    1.25,
    1.00,
    0.89,
]


SIZES = [
    '012',
    '026',
    '048',
]


for size in SIZES:
    for k in K:
        for m in M:
            print(f"size={size}, k={k}, m={m}")
            with open(f'results-{size}-{k}-{m}.txt', mode='wb') as fp:
                for _ in range(500):
                    proc = subprocess.Popen(
                        f'python read_file.py city_data/AISearchFile{size}.txt | anneal/anneal -k={k} -m={m}',
                        shell=True,
                        stdout=subprocess.PIPE,
                    )
                    fp.write(proc.communicate()[0])
