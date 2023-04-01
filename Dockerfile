FROM golang:1.17.6
WORKDIR home-counter

COPY db db
COPY src src
COPY main.go main.go
COPY go.mod go.mod

RUN go get .
CMD ["go", "run", "main.go"]
