FROM alpine

RUN mkdir -p /usr/app/github/config/.
WORKDIR /usr/app/github/.

COPY github-service .
COPY config/config.json ./config/config.json

RUN apk --no-cache add ca-certificates

CMD ["./github-service"]