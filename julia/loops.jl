import Dates

function simulate()
  start = Dates.now()
  iterations = 500_000
  sample_size = 23

  count = 0
  for i in 1:iterations
    data = rand(1:365, 1, sample_size)
    if length(unique(data)) != sample_size
      count += 1
    end
  end
  println("iterations: ", iterations)
  println("sample-size: ", sample_size)
  percent = count / iterations * 100
  println("percent: ", round(percent,digits=2))
  finish = Dates.now()
  diff = (finish - start) |> Dates.value |> x->x/1000
  println("seconds: ", diff)
end

simulate()
