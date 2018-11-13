import os
import hashlib

h = hashlib.sha1()
h.update(open('aco/aco', 'rb').read())
sig = h.hexdigest()[:10]


small_configs = [
    ['-beta=2', '-pg=0.9', '-rho=0.9', '-G=2000'],
    ['-beta=3', '-pg=0.9', '-rho=0.9', '-G=2000'],
    ['-beta=4', '-pg=0.9', '-rho=0.9', '-G=2000'],
    ['-beta=2', '-pg=0.5', '-rho=0.9', '-G=2000'],
    ['-beta=2', '-pg=0.6', '-rho=0.9', '-G=2000'],
    ['-beta=2', '-pg=0.7', '-rho=0.9', '-G=2000'],
    ['-beta=2', '-pg=0.9', '-rho=0.6', '-G=2000'],
    ['-beta=2', '-pg=0.9', '-rho=0.7', '-G=2000'],
    ['-beta=2', '-pg=0.9', '-rho=0.8', '-G=2000'],
    ['-beta=2', '-pg=0.3', '-rho=0.6', '-G=2000'],
    ['-beta=2', '-pg=0.4', '-rho=0.7', '-G=2000'],
    ['-beta=2', '-pg=0.5', '-rho=0.8', '-G=2000'],
]

large_configs = [
    ['-beta=2', '-pg=0.9', '-rho=0.4', '-G=300'],
]



print(sig)
for config in small_configs:
    print(config)
    for size in ['012', '017', '021', '026', '042', '048', '058']:
        print('   ', size)
        os.system('python read_file.py cityfiles/AISearchFile%s.txt | aco/aco' % (size,))
    print()
    print("------------------------------")
    print()
