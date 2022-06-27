FROM golang:1.17 as build-env

ADD go.* /go/src/

WORKDIR /go/src/

RUN go mod download

COPY *.go /go/src/

RUN CGO_ENABLED=0 go build -o main

FROM gcr.io/distroless/static

COPY --from=build-env /go/src/main /app
CMD ["/app"]