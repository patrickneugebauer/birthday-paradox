using System;
using System.Collections.Generic;

class Program {
  static void Main(string[] args) {
    Simulate();
  }
  static void Simulate() {
    DateTime start = DateTime.Now;
    int iterations = 10_000_000;
    int sampleSize = 23;
    int count = 0;
    Random rnd = new Random();

    for (int i=0; i < iterations; i++) {
      int[] data = new int[365];
      for (int l=0; l < sampleSize; l++) {
        int num = rnd.Next(0, 365);
        if (data[num] == 1) {
          count++;
          break;
        } else {
          data[num] = 1;
        }
      }
    }

    Console.WriteLine($"iterations: {iterations}");
    Console.WriteLine($"sample-size: {sampleSize}");
    double percent = ((double) count / iterations) * 100;
    Console.WriteLine($"percent: {Math.Round(percent, 2)}");
    DateTime end = DateTime.Now;
    double diff = (end - start).TotalSeconds;
    Console.WriteLine($"seconds: {Math.Round(diff, 3)}");
  }
}
