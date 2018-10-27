import json
import sys
import threading
import subprocess
from concurrent.futures import ThreadPoolExecutor
from read_file import read

ALPHAS = [0.99993]

SIZES = [
    ('012', 100),
    ('017', 100),
    ('021', 100),
    ('026', 300),
    ('042', 600),
    ('048', 800),
    ('058', 640),
    ('175', 480),
    ('180', 400),
    ('535', 400),
]


def auto_map(tp, bin):
    for size, times in SIZES:
        matrix = json.dumps(read(f'city_data/AISearchFile{size}.txt')).encode('ascii')
        for a in ALPHAS:
            for i in range(1, times+1):
                tp.submit(task, tp, bin, matrix, size, i, times, a)



def task(tp, bin, matrix, size, done, times, alpha, lock=threading.Lock()):
    with lock:
        print(f"[{done}/{times}] RUN  size={size}, alpha={alpha}")

    proc = subprocess.run(
        [f'{bin}/{bin}', f'-alpha={alpha}'],
        capture_output=True,
        input=matrix,
        )

    with lock:
        with open(f'res/{bin}-{size}-{alpha}.txt', mode='ab') as fp:
            fp.write(proc.stdout)
            fp.flush()
        print(f"[{done}/{times}] DONE size={size}, alpha={alpha}")


if __name__ == '__main__':
    auto_map(
        ThreadPoolExecutor(max_workers=4),
        sys.argv[1],
        )
