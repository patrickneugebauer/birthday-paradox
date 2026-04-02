console.log(JSON.stringify(Duktape.env));
const iterations = Duktape.env.args[2];

function simulate() {
  const start = new Date().getTime()
  const sampleSize = 23

  var count = 0
  for (var i=0; i < iterations; i++) {
    const arr = new Array(365)
    for (var j=0; j < sampleSize; j++) {
      var rand = Math.floor(Math.random()*365)
      if (arr[rand]) {
        count++
        break
      } else {
        arr[rand] = 1
      }
    }
  }
  console.log('iterations: ', iterations)
  console.log('sample-size: ', sampleSize)
  var results = (count / iterations * 100).toFixed(2)
  console.log('percent: ', results)
  var end = new Date().getTime()
  var diff = ((end-start) / 1000).toFixed(3)
  console.log('seconds: ', diff)
}

simulate()
