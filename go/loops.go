package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	simulate()
}

func simulate() {
	start := time.Now().UnixNano()
	iterations := 1000000
	sampleSize := 23

	count := 0
	for i := 0; i < iterations; i++ {
		var data [365]int
		for n := 0; n < sampleSize; n++ {
			number := rand.Intn(365)
			if data[number] == 1 {
				count++
				break
			} else {
				data[number] = 1
			}
		}
	}
	fmt.Printf("iterations: %d\n", iterations)
	fmt.Printf("sample-size: %d\n", sampleSize)
	percent := float64(count) / float64(iterations) * 100
	fmt.Printf("percent: %.2f\n", math.Floor(percent*100)/100) // format with printf
	end := time.Now().UnixNano()
	diff := float64(end-start) / 1000 / 1000 / 1000
	fmt.Printf("seconds: %.3f\n", diff)
}
