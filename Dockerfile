# syntax=docker/dockerfile:1

FROM golang:1.18.3-alpine3.16

WORKDIR /

COPY . .

RUN go mod download && go mod verify

RUN go build -o /app ./cmd/serverd/main.go

FROM golang:1.18.3-alpine3.16

WORKDIR /opt

COPY --from=0 /app .

EXPOSE 5000

#CMD [ "/opt/app" ]
