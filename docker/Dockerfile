FROM golang as builder
LABEL authors="pomo"

ENV CGO_ENABLED=1

WORKDIR /jetono-api/src

COPY .. .
RUN mkdir -p /jetono-api/bin
RUN curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh -o install.sh
RUN bash install.sh
COPY ../go.mod ../go.sum /jetono-api/src/

RUN go mod download

RUN go build -o /jetono-api/bin/jetono-api


RUN mkdir -p /store/keys /store/nsc

ENV NKEYS_PATH=/store/keys
ENV NSC_HOME=/store/nsc

RUN cp /jetono-api/bin/jetono-api /usr/local/bin/jetono-api

EXPOSE 1323

CMD ["jetono-api"]
