docker build . -t elixir-build && docker run --rm -t elixir-build 100000

docker build -f Dockerfile.script . -t elixir-script && docker run --rm -t elixir-script 100000
