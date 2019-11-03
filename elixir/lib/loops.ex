defmodule Loops.CLI do

  def main(args) do
    iterations = args
      |> Enum.at(0)
      |> Integer.parse
      |> elem(0)
    iterate(iterations)
  end

  def iterate(iterations) do
    start = Time.utc_now()
    # iterations = 20000
    sample_size = 23
    sample_range = 1..365
    data = Stream.map(1..iterations, fn(_) ->
      Enum.map(1..sample_size, fn(_) ->
        Enum.random(sample_range)
      end)
    end)
    IO.puts "iterations: #{iterations}"
    IO.puts "sample-size: #{sample_size}"
    percent = data
      |> Stream.map(fn(x) -> MapSet.new(x) end)
      |> Stream.filter(fn(x) -> MapSet.size(x) == sample_size end)
      |> Enum.to_list
      |> length
      |> Kernel./(iterations)
      |> Kernel.*(100)
      |> Float.round(2)
    IO.puts "percent: #{percent}"
    # time
    seconds = Time.utc_now()
      |> Time.diff(start, :millisecond)
      |> Kernel./(1000)
      |> Float.round(3)
    IO.puts "seconds: #{seconds}"
  end

end
