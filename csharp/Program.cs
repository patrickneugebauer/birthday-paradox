using System;
using System.Collections.Generic;

class Program {
  static void Main(string[] args) {
    Simulate();
  }
  static void Simulate() {
    DateTime start = DateTime.Now;
    int iterations = 1000000;
    int sampleSize = 23;
    int count = 0;
    Random rnd = new Random();

    for (int i=0; i < iterations; i++) {
      List<int> data = new List<int>();
      for (int l=0; l < sampleSize; l++) {
        int num = rnd.Next(0, 364);
        if (data.Contains(num)) {
          count++;
          break;
        } else {
          data.Add(num);
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
