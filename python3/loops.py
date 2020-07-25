from random import random
import time
import sys
# import array

def simulate():
    start = time.time()
    iterations = int(sys.argv[1])
    sample_size = 23

    count = 0
    for _ in range(iterations):
        # data = [0] * 365
        data = set() # similar speed
        # data = { -1 } # similar speed
        # data = set([]) # similar speed
        # data = [] # 33% slower
        # data = [None for _ in range(365)] # 3x slower
        # data = array.array('H', [0] * 365) # 3.5x slower

        for _ in range(sample_size):
            rand = int(random() * 365)
            # if data[rand] == 1: # for array and list[] solutions
            if rand in data: # for set and x in list solutions
                count += 1
                break
            else:
                # data[rand] = 1 # for array and list[] solutions
                data.add(rand) # for set solution
                # data.append(rand) # for x in list solution

    print("iterations:", iterations)
    print("sample-size:", sample_size)
    results = round(count / iterations * 100, 2)
    print("percent:", results)
    end = time.time()
    diff = round(end-start, 3)
    print("seconds:", diff)

simulate()
