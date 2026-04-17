# variables
iterations = process.argv[2]
start = process.hrtime.bigint()
sampleSize = 23

# generate data
count = 0
for i in [1..iterations]
  arr = new Array(365)
  for j in [1..sampleSize]
    rand = Math.floor(Math.random()*365)
    if arr[rand]
      count++
      break
    else
      arr[rand] = 1

# calcs
results = (count / iterations * 100).toFixed(2)
end = process.hrtime.bigint()
diff = (Number(end-start) / 1_000_000_000).toFixed(6)

# output
console.log("iterations: #{iterations}")
console.log("sample-size: #{sampleSize}")
console.log("percent: #{results}")
console.log("seconds: #{diff}")
