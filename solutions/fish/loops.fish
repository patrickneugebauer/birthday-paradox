set iterations $argv[1]

function simulate
  set sample_size 23
  set start (math (date +%s%N) / 1000 / 1000 / 1000)
  set count 0

  for i in (seq $iterations)
    set data (seq 365)
    for s in (seq $sample_size)
      set num (math --scale=0 (random) \* 364 / 32767 + 1)
      if test $data[$num] -eq 0
        set count (math $count + 1)
        break
      else
        set data[$num] 0
      end
    end
  end

  echo "iterations: $iterations"
  echo "sample-size: $sample_size"
  set percent (math --scale=2 $count \* 100 / $iterations)
  echo "percent: $percent"
  set end (math (date +%s%N) / 1000 / 1000 / 1000)
  set diff (math --scale=3 $end - $start)
  echo "seconds: $diff"
end

simulate
