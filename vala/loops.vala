class HelloWorld : GLib.Object {
  public static int main(string[] args) {
    //  constants
    int64 start = get_real_time();
    const int iterations = 2000000;
    const int sample_size = 23;

    //  generate data
    int count = 0;
    for (int i = 0; i < iterations; i++) {
      int[] sample = new int[365];
      for (int s = 0; s < sample_size; s++) {
        int data = Random.int_range(0, 365);
        if (sample[data] == 1) {
          count++;
          break;
        } else {
          sample[data] = 1;
        }
      }
    }

    //  calcs
    double percent = count * 100.0 / iterations;
    int64 finish = get_real_time();
    int64 microseconds = finish - start;
    double seconds = microseconds / 1000000.0;

    //  output
    stdout.printf("iterations: %d\n", iterations);
    stdout.printf("sample-size: %d\n", sample_size);
    stdout.printf("percent: %.2f\n", percent);
    stdout.printf("seconds: %.3f\n", seconds);

    //  exit code
    return 0;
  }
}
