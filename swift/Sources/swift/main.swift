import Foundation

func simulate() {
  let start = Date()
  let iterations = 100000
  let sampleSize = 23

  var count = 0
  for _ in 1...iterations {
    var data = Array(repeating: 0, count: 365)
    for _ in 1...sampleSize {
      let num = Int.random(in: 0...364)
      if data[num] == 1 {
        count += 1
        break
      } else {
        data[num] = 1
      }
    }
  }
  print("iterations: \(iterations)")
  print("sample-size: \(sampleSize)")
  let percent = Double(count) / Double(iterations) * 100
  print("percent: \( String(format: "%.2f", percent) )")
  let end = Date()
  let diff = end.timeIntervalSince(start)
  print("seconds: \( String(format: "%.3f", diff) )")
}

simulate()
