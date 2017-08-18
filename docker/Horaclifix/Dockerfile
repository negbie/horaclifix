FROM golang:latest
RUN go get -d -v -u github.com/negbie/horaclifix
WORKDIR /go/src/github.com/negbie/horaclifix/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o horaclifix .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/negbie/horaclifix/horaclifix .
CMD ["./horaclifix", "-h"]

