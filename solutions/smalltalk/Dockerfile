FROM nacyot/smalltalk-gnu:apt

WORKDIR /root

ADD loops.st .

ENTRYPOINT ["gst", "loops.st", "-g", "-a"]
