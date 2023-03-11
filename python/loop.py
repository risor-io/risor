def run():
    i = 0
    x = 0
    while i < 10:
        x += i
    return x


def wt():
    x = 0
    while True:
        x += 1
        if x > 10:
            break
    return x


import dis

dis.dis(wt)
