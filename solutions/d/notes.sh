# dmd - commonly used for local dev
docker build -f Dockerfile.dmd . -t d-dmd-build
docker run -it --rm -t d-dmd-build 2500000

# ldc - commonly used for production builds
docker build -f Dockerfile.ldc . -t d-ldc-build
docker run -it --rm -t d-ldc-build 6000000

# gdc - appears to be unmaintained and broken
# docker build -f Dockerfile.gdc . -t d-gdc-build
# docker run -it --rm -t d-gdc-build 1000000
