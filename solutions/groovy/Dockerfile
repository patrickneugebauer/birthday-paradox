FROM groovy:3.0.8-jre16

# groovy docker image requires `-u root` to access `/root` dir
WORKDIR /home/groovy

COPY loops.groovy .

ENTRYPOINT ["groovy", "loops.groovy"]
