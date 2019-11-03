extern crate rand;
extern crate time;

use rand::Rng;

fn main() {
  let start = time::get_time();
  let args: Vec<String> = std::env::args().collect();
  let iterations = args[1].parse::<i32>().unwrap();
  let sample_size = 23;
  let mut random = rand::thread_rng();

  let mut count = 0;
  for _ in 0..iterations {
    let mut data = [0; 365];
    for _ in 0..sample_size {
      let num = random.gen_range(0, 365);
      if data[num] == 1 {
        count += 1;
        break;
      } else {
        data[num] = 1
      }
    }
  }
  println!("iterations: {}", iterations);
  println!("sample-size: {}", sample_size);
  println!("percent: {:.2}", count as f64 / iterations as f64 * 100.0);
  let end = time::get_time();
  let end_time = end.sec as f64 + end.nsec as f64 / 1_000_000_000.0;
  let start_time = start.sec as f64 + start.nsec as f64 / 1_000_000_000.0;
  let diff = end_time - start_time;
  println!("seconds: {:.3}", diff);
}
