FROM mor1/coffeescript

WORKDIR /root

ADD loops.coffee .

RUN coffee -c loops.coffee

ENTRYPOINT ["node", "loops.js"]
