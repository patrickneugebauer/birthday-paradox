using System;
using System.Linq;
using System.Threading.Tasks;

class ProgramParallel
{
  static void MainParallel(string[] args)
  {
    // params
    DateTime start = DateTime.Now;
    int threads = Int32.Parse(args[1]);
    int iterations = Int32.Parse(args[0]); // cap at n=physical_cores
    const int sampleSize = 23;

    // create and run tasks
    Random rnd = new Random();
    var range = Enumerable.Range(0, threads);
    // alternative create + start task: Task.run(fn)
    var tasks = range.Select(i => new Task<int>(
      () => Simulate(iterations, sampleSize, rnd.Next())) // pass in a random value to seed each thread
    ).ToArray();
    Array.ForEach(tasks, task => task.Start());
    Task.WaitAll(tasks);

    // perform calcs
    var duplicates = tasks.Select(task => task.Result).Sum();
    double percent = ((double)duplicates / iterations / threads) * 100;
    DateTime end = DateTime.Now;
    double diff = (end - start).TotalSeconds;
    var ips = iterations * threads / diff;

    // output
    Console.WriteLine($"threads: {threads}");
    Console.WriteLine($"iterations: {iterations * threads}");
    Console.WriteLine($"sample-size: {sampleSize}");
    Console.WriteLine($"percent: {Math.Round(percent, 2)}");
    Console.WriteLine($"seconds: {Math.Round(diff, 3)}");
    Console.WriteLine($"ips: {String.Format("{0:n0}", Math.Round(ips))}");
  }
  static int Simulate(int iterations, int sampleSize, int seed)
  {
    int count = 0;
    // Random number generators are seeded by time
    // explicitly pass in a random seed to ensure uniqueness across threads
    Random rnd = new Random(seed);

    for (int i = 0; i < iterations; i++)
    {
      int[] data = new int[365];
      for (int l = 0; l < sampleSize; l++)
      {
        int num = rnd.Next(0, 365);
        if (data[num] == 1)
        {
          count++;
          break;
        }
        else
        {
          data[num] = 1;
        }
      }
    }

    return count;
  }
}
