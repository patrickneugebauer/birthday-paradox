FROM python:slim

RUN pip install --no-cache-dir xonsh

WORKDIR /app

COPY loops.xsh .

ENTRYPOINT ["xonsh", "loops.xsh"]
