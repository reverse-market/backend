FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/google/wire/cmd/wire
COPY cmd cmd
COPY pkg pkg
RUN cd cmd/reverse-market && wire
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/reverse-market

FROM alpine

LABEL maintainer="Sergey Kozhin <kozhinsrgeyv@gmail.com>"
WORKDIR /app

COPY --from=builder /build/app .
EXPOSE 8080
CMD ["./app"]
