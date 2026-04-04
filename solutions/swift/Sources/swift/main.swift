import Foundation

func simulate() {
  let start = Date()
  let iterationsOptional = Int(CommandLine.arguments[1])
  guard let iterations = iterationsOptional else {
    print("missing iterations")
    return
  }
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

  // calcs
  let percent = Double(count) / Double(iterations) * 100
  let end = Date()
  let diff = end.timeIntervalSince(start)

  // output
  print("iterations: \(iterations)")
  print("sample-size: \(sampleSize)")
  print("percent: \( String(format: "%.2f", percent) )")
  print("seconds: \( String(format: "%.6f", diff) )")
}

simulate()
