FROM golang:1.13.5

WORKDIR /go/src/github.com/dropoutlabs/privacyai

COPY go.mod go.sum Makefile ./
RUN make bootstrap

COPY . ./
RUN make

FROM golang:1.13.5

WORKDIR /root/

COPY --from=0 /go/src/github.com/dropoutlabs/privacyai/bin/privacy .

CMD ["./privacy"]