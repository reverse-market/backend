FROM golang:alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/google/wire/cmd/wire@latest
COPY cmd cmd
COPY pkg pkg
RUN cd cmd/reverse-market && wire
RUN go build -o app ./cmd/reverse-market

FROM alpine

LABEL maintainer="Sergey Kozhin <kozhinsrgeyv@gmail.com>"
WORKDIR /app

COPY --from=builder /build/app .
EXPOSE 8080
CMD ["./app"]
