FROM golang:1.20 AS build
ADD . /go/src/app
WORKDIR /go/src/app
RUN CGO_ENABLED=0 go build -o api .

FROM alpine:3.12.0 AS runtime
WORKDIR /api
COPY --from=build /go/src/app/api /api
EXPOSE 4646/tcp
ENTRYPOINT ["./api"]
