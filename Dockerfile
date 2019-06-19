FROM golang:1.12

WORKDIR /go/src/github.com/swrap/rdma-ds

COPY . .

RUN ( cd v1/server && go build -o app && mv app ../../ )
RUN ls -ltr
CMD ["./app"]
