FROM golang as builder
LABEL authors="pomo"

ENV CGO_ENABLED=0

COPY .. /jetono-api/src
RUN mkdir -p /jetono-api/bin

COPY ../go.mod /jetono-api/src
COPY ../go.sum /jetono-api/src

RUN cd /jetono-api/src && \
    go mod download

RUN cd /jetono-api/src && \
    go build -o /jetono-api/bin/jetono-api

FROM alpine:3.11

RUN apk add --no-cache ca-certificates

RUN mkdir -p /store/keys /store/nsc

ENV NKEYS_PATH=/store/keys
ENV NSC_HOME=/store/nsc

COPY --from=builder /jetono-api/bin/jetono-api /usr/local/bin/jetono-api

EXPOSE 1323

CMD ["jetono-api"]
