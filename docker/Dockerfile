FROM golang as builder
LABEL authors="pomo"

WORKDIR /jetono-api/src
RUN go install github.com/nats-io/nsc/v2@latest

RUN mkdir -p /jetono-api/bin
COPY ./go.mod ./go.sum /

RUN go mod download

COPY . .

RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build go build -o /jetono-api/bin/jetono-api

RUN mkdir -p /store/keys /store/nsc

ENV NKEYS_PATH=/store/keys
ENV NSC_HOME=/store/nsc

RUN cp /jetono-api/bin/jetono-api /usr/local/bin/jetono-api

EXPOSE 1323

CMD ["jetono-api"]
