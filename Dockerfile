FROM golang:1.16-buster as builder
WORKDIR /app

COPY go.mod . 
COPY go.sum .
RUN go mod download 

COPY twitterBot.go .
RUN go build -o ./out/twitterBot

FROM gcr.io/distroless/base
COPY --from=builder /app/out/twitterBot /twitterBot
CMD ["/twitterBot"]
