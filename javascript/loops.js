function simulate() {
  var start = new Date().getTime()
  var iterations = 1000000
  var sampleSize = 23

  var count = 0
  for (var i=0; i < iterations; i++) {
    var arr = new Array(365)
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
  console.log(`iterations: ${iterations}`)
  console.log(`sample-size: ${sampleSize}`)
  var results = (count / iterations * 100).toFixed(2)
  console.log(`percent: ${results}`)
  var end = new Date().getTime()
  var diff = ((end-start) / 1000).toFixed(3)
  console.log(`seconds: ${diff}`)
}

simulate()
