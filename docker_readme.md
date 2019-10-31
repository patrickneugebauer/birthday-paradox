# What

Docker is used throughout this repository to create/run the many different languages.

## Usage

Run some of these few commands inside each language directory in order to build and run code.

- Docker command for building image (verbose output)
    - `docker build --no-cache .`
- Docker command per language to build and run
    - `docker run --rm $(docker build --no-cache --quiet .)`
- Docker command to terminal into a built image
    - `docker run -it --entrypoint=/bin/bash <[intermediate]image-id>`

## Creating New Dockerfiles

Copy and paste this into a new `Dockerfile` to get started:

```dockerfile
FROM ubuntu # Base Ubuntu image

WORKDIR /root # working directory for project

RUN apt update && apt install -y <whatever package> # Run any commands for the image

COPY loops.<file ext> . # Copy project files into image

RUN <build project> # if necessary

ENTRYPOINT ["./loops"] # execute command to run project
```
