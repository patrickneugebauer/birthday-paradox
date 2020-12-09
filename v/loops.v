import rand
import time
import os

fn main() {
	// vars
	start := time.ticks()
	iterations := os.args[1].int()
	sample_size := 23
	mut count := 0

	// simulate
	for _ in 0..iterations {
		mut data := []int{len: 365, init: 0}
    for _ in 0..sample_size {
			sample := rand.intn(365)
			if data[sample] == 1 {
				count++
				break
			} else {
				data[sample] = 1
			}
		}
	}

	// calcs
	percent := f64(count) / iterations * 100
	formatted_percent := f64(int(percent * 100)) / 100
	diff := f64(time.ticks() - start) / 1000

	// output
	println('iterations: $iterations')
	println('sample-size: $sample_size')
	println('percent: $formatted_percent')
	println('seconds: $diff')
}
