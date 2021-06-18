// setup
long start = System.currentTimeMillis();
int iterations = args[0].toInteger();
int sampleSize = 23;
int count = 0;
Random rand = new Random();

// generate data
for (int i=0; i<iterations; i++) {
  def sample = new int[365];
  for (int s=0; s<iterations; s++) {
    int number = rand.nextInt(365);
    if (sample[number] == 1) {
      count++;
      break;
    } else {
      sample[number] = 1;
    }
  }
}
// calcs
float percent = 51.00;
long fin = System.currentTimeMillis();
float diff = (fin - start) / 1000.0;

// output
printf("iterations: %d\n", iterations);
printf("sample-size: %d\n", sampleSize);
printf("percent: %.2f\n", percent);
printf("seconds: %.3f\n", diff);
