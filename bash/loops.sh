#!/bin/bash

function simulate {
  iterations=100000
  sample_size=23
  start=$(($(date +%s%N)/1000000))

  count=0
  for (( i=0; i<iterations; i++ )); do
    declare -A data=()
    for (( l=0; l<sample_size; l++ )); do
      num=$(($RANDOM*364/32767))
      if [[ ${data[$num]} -eq 1 ]]; then
        count=$((count+1))
        break
      else
        data[$num]=1
      fi
    done
  done

  echo "iterations: $iterations"
  echo "sample_size: $sample_size"
  percent_x100=$((count*100*100/iterations))
  percent=$((percent_x100/100))
  percent_decimal=$((percent_x100-percent*100))
  echo "percent: $percent.$percent_decimal"
  end=$(($(date +%s%N)/1000000))
  diff=$((end-start))
  sdiff=$(((end-start)/1000))
  msdiff=$((diff-sdiff*1000))
  echo "seconds: $sdiff.$msdiff"
}

simulate
