import os
import hashlib

h = hashlib.sha1()
h.update(open('aco/aco', 'rb').read())
sig = h.hexdigest()[:10]


configs = {
    "012": [[]],
    "017": [[]],
    "021": [[]],
    "026": [[]],
    "042": [[]],
    "048": [
        ['-beta=3', '-pg=0.4', '-rho=0.9'],
        ['-beta=3', '-pg=0.5', '-rho=0.9'],
        ['-beta=3', '-pg=0.6', '-rho=0.9'],
        ['-beta=3', '-pg=0.7', '-rho=0.9', '-G=5000'],
        ['-beta=2', '-pg=0.7', '-rho=0.9', '-G=5000'], # best
        ['-beta=1', '-pg=0.7', '-rho=0.9', '-G=5000'],
        ['-beta=1', '-pg=0.7', '-rho=0.9', '-G=5000'],
        ['-beta=3', '-pg=0.8', '-rho=0.9'],
        ],
    "058": [[]],
    "175": [['-beta=3', '-pg=0.9', '-rho=0.7', '-G=500']],
    "180": [['-beta=3', '-pg=0.9', '-rho=0.7', '-G=500']],
    "535": [['-beta=2', '-pg=0.9', '-rho=0.4', '-G=300']],
}


print(sig)
for size, configs in configs.items():
    for config in configs:
        file = "aco_results/aco-{hash}-{size}-{args}.txt".format(
                hash=sig,
                args=''.join(config),
                size=size,
                )
        if os.path.exists(file):
            print(" - " + size + " " + ' '.join(config or ['defaults']) + " (skip)")
            continue
        print(" - " + size + " " + ' '.join(config or ['defaults']))
        cmd = (
            'python read_file.py cityfiles/AISearchFile{size}.txt'
            ' | aco/aco {args}'
            ' 2> /dev/null > {file}'
            )
        os.system(cmd.format(
            size=size,
            args=' '.join(config),
            file=file,
            ))
