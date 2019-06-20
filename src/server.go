package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const version = "0.1"
const message = `AvailableEndpoints:
	/getpfs - returns list of available vfs
	
	Version: %s
	`

func getNodePFs() ([]*PF, error) {
	pfs := make([]*PF, 2)
	for i := range pfs {
		vf := &VF{
			MaxTxRate: 0,
			MinTxRate: 0,
		}
		vfs := []*VF{vf}
		pfs[i] = &PF{
			UsedTxRate:     100,
			CapacityTxRate: 100,
			VFs:            vfs,
		}
	}
	return pfs, nil
}

func getNodeData() ([]byte, error) {
	pfs, err := getNodePFs()
	if err != nil {
		return nil, fmt.Errorf("Failed to get Node data: %s", err)
	}
	data, err := json.Marshal(pfs)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal node data: %s", err)
	}
	return data, nil
}

func CreateServer(port string) *http.Server {
	http.HandleFunc("/getpfs", func(w http.ResponseWriter, r *http.Request) {
		data, err := getNodeData()
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s", err)
		} else {
			fmt.Fprintf(w, "%s", data)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, message, version)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: nil}

	log.Println("Starting Server on port: " + port)
	return server
}
