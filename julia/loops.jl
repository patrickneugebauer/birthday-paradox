import Dates

function simulate()
  start = Dates.now()
  iterations = 100_000
  sample_size = 23

  count = 0
  for i in 1:iterations
    data = []
    for x in 1:sample_size
      num = rand(1:365)
      if num in data
        count+=1
        break
      else
        push!(data, num)
      end
    end
  end
  println("iterations: ", iterations)
  println("sample-size: ", sample_size)
  percent = count / iterations
  println("percent: ", round(percent,digits=2))
  finish = Dates.now()
  diff = (finish - start) |> Dates.value |> x->x/1000
  println("seconds: ", diff)
end

simulate()
