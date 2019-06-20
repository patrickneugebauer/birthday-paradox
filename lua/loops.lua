-- constants
start = os.clock()
iterations = 200000
sample_size = 23
count = 0

-- generate data
for i=1,iterations do
  data = {}
  for s=1,sample_size do
    sample = math.random(0, 365)
    if data[sample] == 1 then
      count = count + 1
      break
    else
      data[sample] = 1
    end
  end
end

-- calcs
percent = count / iterations * 100
finish = os.clock()
seconds = finish - start

-- output
print(string.format("iterations: %d", iterations))
print(string.format("sample-size: %d", sample_size))
print(string.format("percent: %.2f", percent))
print(string.format("seconds: %.3f", seconds))
