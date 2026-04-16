if [ -z "$1" ]; then
    echo "Error: enter a decription."
    exit 1
fi
dir="plans"
ts="$(date -u +"%Y%m%dT%H%M%S")"
desc="$1"
fname="$dir/$ts-$desc.md"
touch "$fname"
echo "$fname"
