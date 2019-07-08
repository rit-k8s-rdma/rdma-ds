package main

import (
	"fmt"
	"log"
	"time"

	"github.com/swrap/sriovnet"
)

func configSriov(pfNetdevName string) error {
	var err error

	//enables SRIOV for all devices with cur number of vfs = 0 and has sriov
	log.Println("INFO: enabling sriov (if not already enabled) on netdevice:", pfNetdevName)
	err = sriovnet.EnableSriov(pfNetdevName)
	if err != nil {
		return fmt.Errorf("failed to enable sriov for netdevice=%s: %s", pfNetdevName, err)
	}

	log.Println("INFO: getting pf netdevice handle on:", pfNetdevName)
	pfHandle, err := sriovnet.GetPfNetdevHandle(pfNetdevName)
	if err != nil {
		return fmt.Errorf("failed to get PF handle for netdevice=%s: %s", pfNetdevName, err)
	}

	log.Println("INFO: configuring vfs on netdevice:", pfNetdevName)
	err = sriovnet.ConfigVfs(pfHandle, true)
	if err != nil {
		return fmt.Errorf("failed to config netdevice=%s: %s", pfNetdevName, err)
	}
	return nil
}

func main() {
	const sleepTime = 30 * 1000 * time.Millisecond

	log.Printf("INFO: Starting up...\n")
	//need to wait some time for kernel module to load, else it will error out with larger sriov numbers
	time.Sleep(sleepTime)

	log.Printf("INFO: Gathering system information...\n")
	devices := sriovnet.GetAllRdmaSriovSupportedDevices()
	if len(devices) == 0 {
		log.Println("ERROR: no SRIOV supported PF devices found on your system!")
	}
	for _, device := range devices {
		for i := 0; i < 3; i++ {
			log.Printf("INFO: Attempt %d, configuring SRIOV on netdevice=%s, number of devices=%d\n", i, device, len(devices))
			err := configSriov(device)
			if err != nil {
				log.Println("ERROR: Failed to configure sriov: ", err)
				time.Sleep(sleepTime)
			} else {
				//SUCCESS!
				break
			}
		}
	}
}
