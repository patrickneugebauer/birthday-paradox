#!/bin/sh
filename=$1
iterations=$2
# get fullpath over relative to use in subshell
# we need to execute in a subshell because filename is fixed in scheme source
fullpath="$(realpath $(dirname $0))"

# put arg into file to read
readfile="scheme-input.txt"
command="cd $fullpath \
  && echo $2 > $readfile \
  && scheme --quiet < $1"
# display command before executing
echo "\"\$($command)"\"
# display result of subshell
echo "$(eval $command)"
