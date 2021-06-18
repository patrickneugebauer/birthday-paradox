// import java.util.Random
import kotlin.random.*
import kotlin.system.*
import kotlin.math.*

fun main(args: Array<String>) {
  val iterations = args[0].toInt()
  simulate(iterations)
}

fun simulate(iterations: Int) {
  // val start = System.currentTimeMillis()
  val start = getTimeMillis()
  val sample_size = 23

  var count = 0
  for (i in 1..iterations) {
    val data = IntArray(365)
    for (n in 1..sample_size) {
      val number = Random.nextInt(0,365)
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
  val formattedPercent = round(percent * 100) / 100
  println("percent: ${formattedPercent}")
  val end = getTimeMillis()
  val diff = (end - start) / 1000.0
  println("seconds: ${diff}")
}
