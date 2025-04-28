FROM golang:1.24.2-alpine AS builder
WORKDIR /app

COPY ./main.go ./
COPY ./go.mod  ./
COPY ./api ./api

RUN go get
RUN go build -o /server

FROM gcr.io/distroless/static-debian12:nonroot-38005a5a303edb04eed33167524f55fbdfafe8b9
WORKDIR /root/
COPY --from=builder /server ./app/server
CMD ["./app/server"]