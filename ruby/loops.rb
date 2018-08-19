def simulate
  start = Time.new
  iterations = 100 * 1000
  sample_size = 23

  count = 0
  iterations.times do |l|
    data = []

    sample_size.times do |i|
      num = rand(365)
      if data.include? num
        count += 1
        break
      else
        data[i] = num
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
