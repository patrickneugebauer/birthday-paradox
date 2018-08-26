#include <stdio.h>
#include <time.h>
#include <stdlib.h>
#include <stdbool.h>
#include <sys/time.h>

#define iterations 10000000
#define sampleSize 23

void simulate();

int main() {
  simulate();
}

void simulate() {
  struct timeval start, end;
  gettimeofday(&start, NULL);
  srand(start.tv_usec);

  int duplicates = 0;
  for (int x=0; x < iterations; x++) {
    int data[365] = {};
    for (int i=0; i < sampleSize; i++) {
      int number = (float) rand() / RAND_MAX * 365;
      // check for number in array
      // bool inArray = false;
      // for (int c=0; c<i; c++) {
      //   if (data[c] == number) inArray = true;
      // }
      if (data[number] == 1) {
        // printf("dupe\n");
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
