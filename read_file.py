def parse_int(x):
    b = ""
    for char in x:
        if char in '0123456789':
            b += char
    if not b:
        return
    return int(b)


def read(filename):
    x = open(filename).read()
    x = ''.join(x.split())
    data = x.split(',')

    name = data[0].split('=')[1]
    size = data[1].split('=')[1]
    size = int(size)

    lengths = [parse_int(x) for x in data[2:]]
    lengths = [x for x in lengths if x is not None]

    M = [[0] * size for _ in range(size)]
    k = 0
    for i in range(size):
        for j in range(i + 1, size):
            M[i][j] = lengths[k]
            M[j][i] = lengths[k]
            k += 1
    return name, M
