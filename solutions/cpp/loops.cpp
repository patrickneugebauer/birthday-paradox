#include <stdio.h>
#include <time.h>
#include <stdlib.h>
#include <stdbool.h>
#include <sys/time.h>

#define sampleSize 23

void simulate(int interations);

int main(int argc, char *argv[]) {
  int iterations = atoi(argv[1]);
  simulate(iterations);
  return 0;
}

void simulate(int iterations) {
  struct timeval start, end;
  gettimeofday(&start, NULL);
  srand(start.tv_usec);

  int duplicates = 0;
  for (int x=0; x < iterations; x++) {
    int data[365] = {};
    for (int i=0; i < sampleSize; i++) {
      int number = (float) rand() / RAND_MAX * 365;
      if (data[number] == 1) {
        duplicates++;
        break;
      } else {
        data[number] = 1;
      }
    }
  }
  printf("iterations: %d\n", iterations);
  printf("sample-size: %d\n", sampleSize);
  double results = duplicates * 100.0 / iterations;
  printf("percent: %0.2f\n", results);
  gettimeofday(&end, NULL);
  double startTime = start.tv_sec + start.tv_usec / 1000000.0;
  double endTime = end.tv_sec + end.tv_usec / 1000000.0;
  float diff = endTime - startTime;
  printf("seconds: %.3f\n", diff);
}
