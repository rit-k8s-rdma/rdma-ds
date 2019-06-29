package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/swrap/rdma-ds/src"
)

const (
	configFilePath = "/k8s-rdma-sriov/config.json"
)

func main() {
	//load in config
	log.Printf("Reading in config file from: %s\n", configFilePath)
	configFileData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading in config file[%s]: %s\n", configFilePath, err)
	}

	//marshal config
	var config src.UserConfig
	if err := json.Unmarshal(configFileData, &config); err != nil {
		log.Fatalf("Error unmarshalling config file: %s\n", err)
	}

	//print out config
	output, err := json.MarshalIndent(&config, " ", "\t")
	if err != nil {
		log.Fatalf("Error marshalling config to print out: %s\n", err)
	}
	log.Printf("Finished Loading Config:\n%s\n", string(output))

	if len(config.PfNetdevices) == 0 {
		log.Println("ERROR: no PfNetdevices set in configuration, this node cannot schedule any RDMA pods")
	}

	if len(config.PfMaxBandwidth) == 0 {
		log.Println("ERROR: no PfMaxBandwidth set in configuration, this node cannot schedule any RDMA pods")
	}

	if len(config.PfMaxBandwidth) != len(config.PfNetdevices) {
		log.Println("ERROR: len(config.PfMaxBandwidth) != len(config.PfNetdevices), this node cannot schedule any RDMA pods")
	}

	//start up server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "54005"
	}
	server := src.CreateServer(port, config)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %s\n", err)
	}
}
