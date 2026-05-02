const iterations = process.argv[2];

function simulate() {
  // declarations
  const start = process.hrtime.bigint();
  const sampleSize = 23
  let count = 0
  // loop
  for (let i = 0; i < iterations; i++) {
    const list = new Int16Array(sampleSize);
    for (let j = 0; j < sampleSize; j++) {
      const rand = Math.floor(Math.random() * 365)
      if (list.includes(rand)) {
        count++
        break
      } else {
        list[j] = rand
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
