package v1_test

import (
	"context"
	"testing"
	"time"

	"net/http"

	"github.com/swrap/rdma-ds/v1"
)

func Test_getPfs(t *testing.T) {
	host := "localhost"
	port := "45005"
	var server *http.Server

	go func() {
		server = v1.CreateServer(port)
		if err := server.ListenAndServe(); err != nil {
			t.Fatalf("Failed to run server: %s", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	_, err := v1.GetNodeInfo(host, port)
	if err != nil {
		t.Error("Failed to get node info")
	}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		t.Fatal("Failed shutdown node")
	}
}
