import os
import random
import strutils
import times
import math

# variables
let start = epochTime()
let iterations = parseInt(paramStr(1))
const sampleSize = 23
var count = 0

proc roundTo(number: float, places: float): float =
  let shifter = pow(10, places)
  int(number * shifter) / int(shifter)

# setup
randomize()
for i in countup(1, iterations):
  type Sample = array[1..365, bool]
  var sample: Sample
  for s in countup(1, sampleSize):
    let day = random(1..365)
    if sample[day]:
      inc(count)
      break
    else:
      sample[day] = true

# calcs
let percent = count / iterations * 100
let fin = epochTime()
let diff = fin - start

# output
echo "iterations: " & $iterations
echo "sample-size: " & $sampleSize
echo "percent: " & $roundTo(percent, 2)
echo "seconds: " & $roundTo(diff, 3)
