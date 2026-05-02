const iterations = process.argv[2];

function simulate() {
  // declarations
  const start = process.hrtime.bigint();
  const sampleSize = 23
  let count = 0
  const arr = []
  let rand
  // loop
  for (let i = 0; i < iterations; i++) {
    for (let j = 0; j < sampleSize; j++) {
      rand = Math.floor(Math.random() * 365)
      if (arr[rand] === i) {
        count++
        break
      } else {
        arr[rand] = i
      }
    }
  }
  // calcs
  const results = (count / iterations * 100).toFixed(2)
  const end = process.hrtime.bigint();
  const diff = (Number(end - start) / 1_000_000_000).toFixed(6)
  // log
  console.log(`iterations: ${iterations}`)
  console.log(`sample-size: ${sampleSize}`)
  console.log(`percent: ${results}`)
  console.log(`seconds: ${diff}`)
}

simulate()
