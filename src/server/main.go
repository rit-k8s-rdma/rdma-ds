package main

import (
	"log"
	"os"

	"github.com/swrap/rdma-ds/src"
	"github.com/swrap/sriovnet"
)

func main() {
	//load in config
	log.Printf("INFO: Gathering system information...\n")

	//Looks up system information to SRIOV enabled Devices
	devices := sriovnet.GetAllRdmaSriovSupportedDevices()
	if len(devices) == 0 {
		log.Println("ERROR: no SRIOV enabled PF devices found on your system!")
	}
	config := src.SystemConfig{
		PfNetDevices: make([]src.PfNetDevice, 0),
	}
	for _, device := range devices {
		rate, err := sriovnet.GetPfMaxSendingRate(device)
		if err != nil {
			log.Printf("ERROR: PF device[%s] max sending rate not found, setting to 0: %s\n", device, err)
			continue
		}
		log.Printf("INFO: Found SRIOV enabled PF device[%s] with rate in bytes [%d]\n", device, rate)
		config.PfNetDevices = append(config.PfNetDevices, src.PfNetDevice{
			Name:           device,
			MaxSendingRate: rate,
		})
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
