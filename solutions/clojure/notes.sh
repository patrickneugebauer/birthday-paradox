docker build . -t clojure-build
docker run --rm -t clojure-build 100000

docker build -f Dockerfile.parallel . -t clojure-build-parallel
docker run --rm -t clojure-build-parallel
