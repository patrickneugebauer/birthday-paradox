object Loops {
  def main(args: Array[String]) = {
    // constants
    val start = System.currentTimeMillis()
    val iterations = 500000
    val sampleSize = 23

    // create data
    val random = scala.util.Random
    val createDataPoint = (_: Int) => random.nextInt(365)
    val createSample = (_: Int) => (1 to sampleSize).map(createDataPoint)
    val data = (1 to iterations).map(createSample)

    // calcs
    val count = data.filter(sample => (sample.distinct.length != sampleSize)).length
    val percent = count.toFloat / iterations * 100
    val finish = System.currentTimeMillis()
    val seconds = (finish - start).toFloat / 1000

    // output
    println(s"iterations: $iterations")
    println(s"sample-size: $sampleSize")
    println(f"percent: $percent%.2f")
    println(f"seconds: $seconds%.3f")
  }
}
