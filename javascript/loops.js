const iterations = process.argv[2];

function simulate() {
  const start = new Date().getTime()
  const sampleSize = 23

  let count = 0
  for (let i=0; i < iterations; i++) {
    const arr = new Array(365)
    for (let j=0; j < sampleSize; j++) {
      const rand = Math.floor(Math.random()*365)
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
  const end = new Date().getTime()
  const diff = ((end-start) / 1000).toFixed(3)
  console.log(`seconds: ${diff}`)
}

simulate()
