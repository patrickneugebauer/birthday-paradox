FROM node

WORKDIR /root

# do this first to avoid regenerating these steps when code changes
RUN npm i -g esy
ADD package.json .
ADD reason.opam .
RUN esy

# these files will change
ADD dune-project .
ADD bin bin
ADD lib lib

ENTRYPOINT ["esy", "x", "Loops"]
