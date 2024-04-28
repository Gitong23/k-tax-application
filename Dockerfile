FROM golang:1.22.1-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/k-taxes .

# ====================

FROM alpine:3.16.2
COPY --from=build-base /app/out/k-taxes /app/k-taxes

CMD ["/app/k-taxes"]