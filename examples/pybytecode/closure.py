from dis import dis


def incrementer():
    x = 0

    def increment():
        nonlocal x
        x += 1
        return x

    return increment


inc = incrementer()
print(inc())
print(inc())

print("dis for inc:")
dis(inc)

print("----")
print("dis for incrementer:")
dis(incrementer)
