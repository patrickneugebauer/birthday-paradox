#!/bin/bash

function padl {
  string=$1
  result=$string
  length=$2
  char=$3
  original_length=${#string}
  diff=$(($length-$original_length))
  for ((i=1; i<=$diff; i++)); do
    result="$char$result"
  done
  echo $result
}

function simulate {
  iterations=5000
  sample_size=23
  start=$(($(date +%s%N)/1000/1000))

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
  echo "sample-size: $sample_size"
  percent_x100=$((count*100*100/iterations))
  percent=$((percent_x100/100))
  percent_decimal=$((percent_x100-percent*100))
  echo "percent: $percent.$percent_decimal"
  end=$(($(date +%s%N)/1000/1000))
  diff=$((end-start))
  sdiff=$(((end-start)/1000))
  msdiff=$((diff-sdiff*1000))
  msdiffstring=$(padl $msdiff 3 0)
  echo "seconds: $sdiff.$msdiffstring"
}

simulate
