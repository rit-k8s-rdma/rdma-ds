FROM golang:1.12 as builder

WORKDIR /go/src/github.com/swrap/rdma-ds

COPY bin/init ./app

# RUN ( cd src/server && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && mv app ../../ )

FROM scratch

WORKDIR /bin
COPY --from=builder /go/src/github.com/swrap/rdma-ds/app .

CMD ["./app"]
