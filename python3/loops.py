from random import randint
import time

def simulate():
    start = time.time()
    iterations = 100 * 1000
    sample_size = 23

    count = 0
    for x in range(iterations):
        data = []

        for i in range(sample_size):
            rand = randint(1, 365)
            if rand in data:
                count += 1;
                break
            else:
                data.append(rand)

    print("interations:", iterations)
    print("sample-size:", sample_size)
    results = round(count / iterations * 100, 2);
    print("percent:", results)
    end = time.time()
    diff = round(end-start, 3)
    print("seconds:", diff)

simulate()
