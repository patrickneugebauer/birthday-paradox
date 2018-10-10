def simulate
  start = Time.new
  iterations = 200_000
  sample_size = 23

  count = 0
  iterations.times do |l|
    data = Array.new(365)

    sample_size.times do |i|
      num = rand(365)
      if data[num] == 1
        count += 1
        break
      else
        data[num] = 1
      end
    end

  end
  puts "iterations: #{iterations}"
  puts "sample-size: #{sample_size}"
  results = (count.to_f / iterations * 100).round(2)
  puts "percent: #{results}"
  fin = Time.new
  diff = (fin - start).round(3)
  puts "seconds: #{diff}"
end

simulate
