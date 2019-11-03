def simulate
  start = Time.local
  iterations = ARGV.first.to_i
  sample_size = 23

  count = 0
  iterations.times do |l|
    data = [] of Int32

    sample_size.times do |i|
      num = rand(365)
      if data.includes? num
        count += 1
        break
      else
        data << num
      end
    end

  end
  puts "iterations: #{iterations}"
  puts "sample-size: #{sample_size}"
  results = (count.to_f / iterations * 100).round(2)
  puts "percent: #{results}"
  fin = Time.local
  time_diff = (fin - start)
  diff = time_diff.seconds + time_diff.milliseconds / 1000.0
  puts "seconds: #{diff}"
end

simulate
