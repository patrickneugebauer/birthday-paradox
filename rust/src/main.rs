extern crate rand;
extern crate time;

use rand::Rng;

fn main() {
  let start = time::get_time();
  let iterations = 1_000_000;
  let sample_size = 23;
  let mut random = rand::thread_rng();

  let mut count = 0;
  for i in 0..iterations {
    let mut data: Vec<i32> = Vec::new();
    for n in 0..sample_size {
      let num = random.gen_range(0, 365);
      if data.contains(&num) {
        count += 1;
        break;
      } else {
        data.push(num)
      }
    }
  }
  println!("iterations: {}", iterations);
  println!("sample-size: {}", sample_size);
  println!("percent: {:.2}", count as f64 / iterations as f64 * 100.0);
  let end = time::get_time();
  let endTime = end.sec as f64 + end.nsec as f64 / 1_000_000_000.0;
  let startTime = start.sec as f64 + start.nsec as f64 / 1_000_000_000.0;
  let diff = endTime - startTime;
  println!("seconds: {:.3}", diff);
}
