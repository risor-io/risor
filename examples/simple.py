def add(a, b):
    return a + b


def callit():
    return add(1, 2)


def setvar():
    x = 1
    y = 2
    return x + y

def value():
    return 4

def testit():
    { 1: 2 }
    { 1: 2 }
    { 1: 2 }
    return 7

# print(testit())


import dis

dis.dis(testit)
