from random import random
import time
import sys

def simulate():
    start = time.time()
    iterations = int(sys.argv[1])
    sample_size = 23

    count = 0
    for x in range(iterations):
        data = [0] * 365

        for i in range(sample_size):
            rand = int(random() * 365)
            if data[rand] == 1:
                count += 1;
                break
            else:
                data[rand] = 1;

    print "iterations:", iterations
    print "sample-size:", sample_size
    results = round(float(count) / iterations * 100, 2);
    print "percent:", results
    end = time.time()
    diff = round(end-start, 3)
    print "seconds:", diff

simulate()
