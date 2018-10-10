import java.util.Random

fun main(args: Array<String>) {
  simulate()
}

fun simulate() {
  val start = System.currentTimeMillis();
  val iterations = 2000 * 1000
  val sample_size = 23
  val random = Random()

  var count = 0
  for (i in 1..iterations) {
    val data = IntArray(365);
    for (n in 1..sample_size) {
      val number = random.nextInt(365)
      if (data[number] == 1) {
        count++
        break
      } else {
        data[number] = 1
      }
    }
  }
  println("iterations: ${iterations}")
  println("sample-size: ${sample_size}")
  val percent = count * 100.0 / iterations
  println("percent: %.2f".format(percent))
  val end = System.currentTimeMillis();
  val diff = (end - start) / 1000.0
  println("seconds: ${diff}")
}
