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

    seconds = Time.utc_now()
      |> Time.diff(start, :millisecond)
      |> Kernel./(1000)
      |> Float.round(3)

    iterations_per_second = total_iterations
      |> Kernel./(seconds)
      |> Kernel.round

    iterations_per_second_per_worker = iterations_per_second
      |> Kernel./(workers)
      |> Kernel.round

    [
      iterations: total_iterations,
      seconds: seconds,
      workers: workers,
      iterations_per_second: iterations_per_second,
      iterations_per_second_per_worker: iterations_per_second_per_worker
    ]
  end

  def generate_randoms(r, c, rng) do
    Enum.map(1..r, fn(_) ->
      Enum.map(1..c, fn(_) ->
        Enum.random(rng)
      end)
    end)
  end

end

range = 1..12
IO.puts "running for #{Enum.at(range, 0)} to #{Enum.at(range, Enum.count(range) - 1)} workers"
results = Enum.map(range, fn(x) ->
  IO.puts "running #{x} workers"
  Loops.iterate (x)
end)
headers = Enum.at(results, 0) |> Keyword.keys |> Enum.join("\t")
rows = Enum.map(results, fn(x) ->
  x |> Keyword.values |> Enum.join("\t")
end)
all_rows = [ headers | rows ]
file_text = Enum.join(all_rows, "\n")

filename = "elixir/stats.txt"
{ :ok, file } = File.open(filename, [:write])
IO.binwrite(file, file_text)
File.close(file)
IO.puts "results written to #{filename}"
