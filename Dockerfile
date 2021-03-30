FROM golang:1.16.2 as build
WORKDIR /bar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=build /bar .
ENTRYPOINT ["/bar"]