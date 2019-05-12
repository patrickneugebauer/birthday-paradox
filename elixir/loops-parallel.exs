defmodule Loops do
  @iterations 12500
  @sample_size 23
  @sample_range 1..365

  def iterate(workers) do
    start = Time.utc_now()
    Enum.map(1..workers, fn(_) ->
      Task.async(fn() ->
        Loops.generate_randoms(@iterations, @sample_size, @sample_range)
          |> Enum.filter(fn(x) ->
            x |> MapSet.new |> MapSet.size |> Kernel.==(@sample_size)
          end)
          |> length
      end)
    end)
      |> Enum.map(fn (x) -> Task.await(x) end)

    total_iterations = @iterations * workers
    IO.puts "#{total_iterations} iterations"

    seconds = Time.utc_now()
      |> Time.diff(start, :millisecond)
      |> Kernel./(1000)
      |> Float.round(3)
    IO.puts "#{seconds} seconds"

    iterations_per_second = total_iterations
      |> Kernel./(seconds)
      |> Kernel.round
    IO.puts "#{iterations_per_second} iterations/sec"

    iterations_per_second_per_worker = iterations_per_second
      |> Kernel./(workers)
      |> Kernel.round
    IO.puts "#{iterations_per_second_per_worker} iterations/sec/worker"
  end

  def generate_randoms(r, c, rng) do
    Enum.map(1..r, fn(_) ->
      Enum.map(1..c, fn(_) ->
        Enum.random(rng)
      end)
    end)
  end

end

Enum.each(1..8, fn (x) ->
  IO.puts "#{x} workers"
  Loops.iterate(x)
  IO.puts "-------------------------"
end)
