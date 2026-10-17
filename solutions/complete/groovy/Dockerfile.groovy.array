FROM groovy:alpine

# groovy docker image requires `-u root` to access `/root` dir
WORKDIR /home/groovy

COPY loops.groovy .

ENTRYPOINT ["groovy", "loops.groovy"]
