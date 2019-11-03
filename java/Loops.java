import java.util.Random;

public class Loops {

  public static void main(String[] args) {
    final int iterations = Integer.parseInt(args[0]);
    simulate(iterations);
  }

  private static void simulate(int iterations) {
    long start = System.currentTimeMillis();
    final int SAMPLE_SIZE = 23;
    Random rand = new Random();
    rand.setSeed(start);
    int duplicates = 0;
      for (long x=0; x < iterations; x++) {
      int data[] = new int[365];
      for (int i=0; i < SAMPLE_SIZE; i++) {
        int number = rand.nextInt(365);
        if (data[number] == 1) {
          duplicates++;
          break;
        } else {
          data[number] = 1;
        }
      }
    }
    System.out.println("iterations: " + iterations);
    System.out.println("sample-size: " + SAMPLE_SIZE);
    double results = duplicates * 100.0 / iterations;
    System.out.format("percent: %.2f%n", results);
    long end = System.currentTimeMillis();
    double diff = (end-start) / 1_000.0;
    System.out.format("seconds: %.3f%n", diff);
  }

}
