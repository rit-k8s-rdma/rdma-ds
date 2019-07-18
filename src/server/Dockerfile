FROM golang:1.12 as builder

WORKDIR /go/src/github.com/rit-k8s-rdma/rit-k8s-rdma-ds

COPY . .

WORKDIR /go/src/github.com/rit-k8s-rdma/rit-k8s-rdma-ds/src/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app
RUN mv app ../../

FROM scratch

WORKDIR /bin
COPY --from=builder /go/src/github.com/rit-k8s-rdma/rit-k8s-rdma-ds/app .

CMD ["./app"]
