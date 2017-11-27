# sudo docker build -t negbie/horaclifix:latest .

FROM golang:latest as builder
RUN go get -d -v -u github.com/negbie/horaclifix
WORKDIR /go/src/github.com/negbie/horaclifix/
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -installsuffix cgo -o horaclifix .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/negbie/horaclifix/horaclifix .
EXPOSE 4739
CMD ["./horaclifix", "-h"]
