import json
import sys
import threading
import subprocess
from concurrent.futures import ThreadPoolExecutor
from read_file import read

ALPHAS = [
    #0.98,
    #0.99,
    #0.992,
    #0.994,
    #0.996,
    #0.9967,
    #0.998,
    #0.9995,
    #0.9997,
    #0.9999,
    0.99993,
]

SIZES = [
    #('012', 200),
    #('017', 200),
    #('021', 200),
    #('026', 120),
    ('042', 224),
    ('048', 200),
    ('058', 160),
    ('175', 120),
    ('180', 80),
    ('535', 40),
]


def auto_map(tp):
    for size, times in SIZES:
        matrix = json.dumps(read(f'city_data/AISearchFile{size}.txt')).encode('ascii')
        for a in ALPHAS:
            for i in range(1, times+1):
                tp.submit(task, tp, matrix, size, i, times, a)



def task(tp, matrix, size, done, times, alpha, lock=threading.Lock()):
    with lock:
        print(f"[{done}/{times}] RUN  size={size}, alpha={alpha}")

    proc = subprocess.run(
        ['anneal/anneal', f'-alpha={alpha}'],
        capture_output=True,
        input=matrix,
        )

    with lock:
        with open(f'results/newswap-{size}-{alpha}.txt', mode='ab') as fp:
            fp.write(proc.stdout)
            fp.flush()
        print(f"[{done}/{times}] DONE size={size}, alpha={alpha}")


if __name__ == '__main__':
    auto_map(ThreadPoolExecutor(max_workers=4))
