FROM golang:1.6
LABEL version="20170727003"
RUN go get github.com/negbie/siprocket
RUN go get github.com/negbie/horaclifix \
    && go install github.com/negbie/horaclifix
COPY ./entrypoint.sh /
RUN chmod +x /entrypoint.sh
EXPOSE 4739
ENTRYPOINT ["/entrypoint.sh"]
