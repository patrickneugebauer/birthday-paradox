import java.util.List;
import java.util.ArrayList;

public class Loops {

  public static void main(String[] args) {
    simulate();
  }

  private static void simulate() {
    long start = System.currentTimeMillis();
    final long ITERATIONS = 100_000;
    final int SAMPLE_SIZE = 23;

    int duplicates = 0;
      for (long x=0; x < ITERATIONS; x++) {
      List<Integer> data = new ArrayList<>();
      for (int i=0; i < SAMPLE_SIZE; i++) {
        int number = (int) Math.round(Math.random() * 365);
        if (data.contains(number)) {
          duplicates++;
          break;
        } else {
          data.add(number);
        }
      }
    }
    System.out.println("iterations: " + ITERATIONS);
    System.out.println("sample-size: " + SAMPLE_SIZE);
    double results = duplicates * 100.0 / ITERATIONS;
    System.out.format("percent: %.2f%n", results);
    long end = System.currentTimeMillis();
    double diff = (end-start) / 1_000.0;
    System.out.format("seconds: %.3f%n", diff);
  }

}
