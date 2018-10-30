import json
import sys
import threading
import subprocess
from concurrent.futures import ThreadPoolExecutor
from read_file import read

CONFIGS = [
    ['anneal', '-alpha=0.99999'],
    ['anneal-sigmoid', '-epsilon=0.00001'],
]

SIZES = [
    ('017', 100),
    ('021', 100),
    ('026', 100),
    ('042', 100),
    ('048', 100),
    ('058', 100),
    ('175', 100),
    ('180', 100),
    ('535', 100),
]


def auto_map(tp):
    for size, times in SIZES:
        matrix = json.dumps(read(f'city_data/AISearchFile{size}.txt')).encode('ascii')
        for bin, config in CONFIGS:
            for i in range(1, times+1):
                tp.submit(task, bin, config, matrix, size, i, times)



def task(bin, config, matrix, size, done, times, lock=threading.Lock()):
    with lock:
        print(f"[{done}/{times}] RUN  size={size} {bin} {config}")

    proc = subprocess.run(
        [f'{bin}/{bin}', config],
        capture_output=True,
        input=matrix,
        )

    with lock:
        with open(f'grind/{bin}-{size}-{config}.txt', mode='ab') as fp:
            fp.write(proc.stdout)
            fp.flush()
        print(f"[{done}/{times}] DONE size={size} {config}")


if __name__ == '__main__':
    auto_map(ThreadPoolExecutor(max_workers=4))
