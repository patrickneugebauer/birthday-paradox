FROM golang:tip-alpine
WORKDIR /app

COPY loops.go .

RUN go build -o loops loops.go

ENTRYPOINT ["./loops"]
# previously was using golang image (non alpine)
# old image is 873MB without binary, 909MB with
# new image is 264MB without binary, 300MB with
# 36MB taken up by binary in image, 24MB locally
