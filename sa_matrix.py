import json
import sys
import threading
import subprocess
from concurrent.futures import ThreadPoolExecutor
from read_file import read

alpha = [
    0.98,
    0.99,
    0.992,
    0.994,
    0.996,
    0.9967,
    0.998,
    0.9995,
]

SIZES = [
    #('012', 500),
    #('017', 500),
    ('021', 500),
    ('026', 300),
    ('042', 280),
    ('048', 250),
    ('058', 200),
    ('175', 150),
    ('180', 75),
    ('535', 20),
]


def task(size, times, alpha, lock=threading.Lock()):
    with lock:
        print(f"Running size={size}, alpha={alpha}")

    matrix = read(f'city_data/AISearchFile{size}.txt')
    input = json.dumps(matrix).encode('ascii')

    with open(f'results/results-{size}-{alpha}.txt', mode='ab') as fp:
        for i in range(times):
            proc = subprocess.run(
                ['anneal/anneal', f'-alpha={alpha}'],
                capture_output=True,
                input=input,
                )
            fp.write(proc.stdout)

    with lock:
        print(f"Done size={size}, alpha={alpha}")


if __name__ == '__main__':
    tp = ThreadPoolExecutor(max_workers=4)
    tasks = []
    for size, times in SIZES:
        for a in alpha:
            f = tp.submit(task, size, times, a)
            tasks.append(f)
    for future in tasks:
        future.result()
