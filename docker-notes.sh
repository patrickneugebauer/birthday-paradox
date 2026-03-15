# notes on building, running, and getting size of images

# build
docker build . -t go

# run
docker run --rm go 1000000

# get image size
docker images | grep go | rev | cut -d " " -f 1 | rev
# => WARNING: This output is designed for human readability. For machine-readable output, please use --format.

docker image ls --filter "reference=*go*" --format "{{.Size}}"
docker image ls go --format "{{.Size}}"
docker image ls go:latest --format "{{.Size}}"
