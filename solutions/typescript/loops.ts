const iterations = parseInt(process.argv[2]);

function simulate() {
  const start = process.hrtime.bigint();
  const sampleSize = 23

  let count = 0
  for (let i = 0; i < iterations; i++) {
    const arr = new Array(365)
    for (let j = 0; j < sampleSize; j++) {
      const rand = Math.floor(Math.random() * 365)
      if (arr[rand]) {
        count++
        break
      } else {
        arr[rand] = 1
      }
    }
  }
  console.log(`iterations: ${iterations}`)
  console.log(`sample-size: ${sampleSize}`)
  const results = (count / iterations * 100).toFixed(2)
  console.log(`percent: ${results}`)
  const end = process.hrtime.bigint();
  const diff = (Number(end - start) / 1_000_000_000).toFixed(6)
  console.log(`seconds: ${diff}`)
}

simulate()
