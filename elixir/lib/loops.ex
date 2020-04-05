defmodule Loops.CLI do

  def main(args) do
    iterations = args
      |> Enum.at(0)
      |> Integer.parse
      |> elem(0)
    iterate(iterations)
  end

  def duplicates([]) do
    false
  end

  def duplicates([x|xs]) do
    if Enum.member?(xs, x), do: true, else: duplicates(xs)
  end

  def iterate(iterations) do
    # data
    start = Time.utc_now()
    sample_size = 23
    sample_range = 1..365
    data = Stream.map(1..iterations, fn(_) ->
      Enum.map(1..sample_size, fn(_) ->
        Enum.random(sample_range)
      end)
    end)
    percent = data
      |> Stream.filter(&Loops.CLI.duplicates/1)
      |> Enum.to_list
      |> length
      |> Kernel./(iterations)
      |> Kernel.*(100)
      |> Float.round(2)
    seconds = Time.utc_now()
      |> Time.diff(start, :millisecond)
      |> Kernel./(1000)
      |> Float.round(3)
    # output
    IO.puts "iterations: #{iterations}"
    IO.puts "sample-size: #{sample_size}"
    IO.puts "percent: #{percent}"
    IO.puts "seconds: #{seconds}"
  end

end
