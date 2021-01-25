FROM ubuntu

RUN mkdir -p /usr/app/github/config/.
WORKDIR /usr/app/github/.

COPY github-service .
COPY config/config.json ./config/config.json

CMD ["./github-service"]