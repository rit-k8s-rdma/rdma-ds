package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rit-k8s-rdma/rit-k8s-rdma-common/rdma_hardware_info"
)

func GetNodeInfo(ip string, port string) ([]*rdma_hardware_info.PF, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/getpfs", ip, port))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	var pfs []*rdma_hardware_info.PF
	if err := json.Unmarshal(data, &pfs); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal data: %s", err)
	}
	return pfs, nil
}
