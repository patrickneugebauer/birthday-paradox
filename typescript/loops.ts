function simulate() {
  const start = new Date().getTime();
  const iterations = 1_000_000;
  const sampleSize = 23;

  let count = 0;
  for (let i = 0; i < iterations; i++) {
    const arr: number[] = [];
    for (let j = 0; j < sampleSize; j++) {
      const rand = Math.floor(Math.random() * 365);
      if (arr.includes(rand)) {
        count++;
        break;
      } else {
        arr[j] = rand;
      }
    }
  }
  console.log(`iterations: ${iterations}`);
  console.log(`sample-size: ${sampleSize}`);
  var results = (count / iterations * 100).toFixed(2);
  console.log(`percent: ${results}`);
  var end = new Date().getTime();
  var diff = ((end-start) / 1000).toFixed(3);
  console.log(`seconds: ${diff}`);
}

simulate();
