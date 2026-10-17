FROM alpine
WORKDIR /app
RUN apk update && apk add --no-cache luau
COPY loops.luau .

ENTRYPOINT ["luau", "loops.luau"]
