package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	sriov "github.com/Mellanox/sriov-cni/src"
)

const (
	version = "latest"
	message = `AvailableEndpoints:
		/getpfs - returns list of available vfs
		
		Version: %s
	`
)

//isVFAllocated reads the sysfs directory for a vf in a pf and determines
//if the vf is in use. It determines if its in use by looking at whether
//any directories exist in /sys/class/net/<pf-name>/device/virtfn<vf-number>/net
//directory. If all directories in that directory are missing, one can assume
//they are being used in a pod.
func isVfAllocated(pf string, vfi uint) (bool, error) {
	vfDir := fmt.Sprintf("/sys/class/net/%s/device/virtfn%d/net", pf, vfi)
	if _, err := os.Lstat(vfDir); err != nil {
		return false, fmt.Errorf("failed to open the virtfn%d dir of the device %s, vfDir[%s] could not be opened: %s", vfi, pf, vfDir, err)
	}
	infos, err := ioutil.ReadDir(vfDir)
	if err != nil || len(infos) == 0 {
		//assume if there are no directories in this directory than VF is in use/allocated
		return true, nil
	}
	//if one or more directories are found, than the VF is available for use
	return false, nil
}

//setVFConfig is responsible for reading the config file in the sysfs directory
//and setting the vf's information. This function only returns an error if the
//directory does not exist. For each line of data in the configuration file,
//if it can't split the configuration line by a ':' delimeter it skips that line
//prints out a warning, but continues to run.
func setVFConfig(pf string, vfi uint, vf *VF) error {
	vfConfigFilePath := fmt.Sprintf("/sys/class/net/%s/device/sriov/%d/config", pf, vfi)
	vfConfigFile, err := os.Open(vfConfigFilePath)
	if err != nil {
		return fmt.Errorf("could not read in vf config file[%s]: %s", vfConfigFilePath, err)
	}
	vfConfigFileReader := bufio.NewReader(vfConfigFile)
	vfConfigFileScanner := bufio.NewScanner(vfConfigFileReader)
	var tmpVal uint64
	const warningConverting string = "WARNING: failed to convert str to number for vf[%d] on device[%s], not setting configuration[%s], continuing onto rest of config...\n"
	for vfConfigFileScanner.Scan() {
		lineText := vfConfigFileScanner.Text()
		configKeyVal := strings.Split(lineText, " : ")
		if len(configKeyVal) != 2 {
			log.Printf("WARNING: splitting config for vf[%d] on device[%s] failed, "+
				"not setting configuration[%s], continuing onto rest of config...\n",
				vfi, pf, lineText)
			continue
		}
		key := strings.TrimSpace(configKeyVal[0])
		val := strings.TrimSpace(configKeyVal[1])
		//reads config file and sets data based on string matching
		switch key {
		case "VF":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.VFNumber = uint(tmpVal)
		case "MAC":
			vf.MAC = val
		case "VLAN":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.VLAN = uint(tmpVal)
		case "QoS":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.QoS = uint(tmpVal)
		case "VLAN Proto":
			vf.VLanProto = val
		case "SpoofCheck":
			vf.SpoofCheck = val
		case "Trust":
			vf.Trust = val
		case "LinkState":
			vf.LinkState = val
		case "MinTxRate":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.MinTxRate = uint(tmpVal)
		case "MaxTxRate":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.MaxTxRate = uint(tmpVal)
		case "VGT+":
			vf.VGTPlus = val
		case "RateGroup":
			tmpVal, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Printf(warningConverting, vfi, pf, lineText)
				continue
			}
			vf.RateGroup = uint(tmpVal)
		default:
			log.Printf("WARNING: Unknown value for vf[%d] pf[%s] for sys config[%s], lineText[%s] Key[%s] Value[%s]\n",
				vfi, pf, vfConfigFilePath, lineText, key, val)
		}
	}
	return nil
}

//getNodeData reads the systems configuration files for each pf
//and vf.
func getNodeData(userConfig UserConfig) ([]*PF, error) {
	//checks pfs and make sure all are enabled
	//orders pfs by the most number of vfs earlier on
	pfs, err := sriov.GetPFs("", userConfig.PfNetdevices)
	if err != nil {
		return nil, err
	}

	nodePfs := make([]*PF, 0)

	//go through every pf and get all pf information
	for _, pf := range pfs {
		totalVfs, err := sriov.GetsriovNumfs(pf)
		if err != nil {
			log.Printf("Error counting vfs for pf[%s]: %s\n", pf, err)
			continue
		}
		tmpNodePf := PF{
			VFs: make([]*VF, 0),
		}
		//get information about each vf that is part of pf
		for ivf := 0; ivf < totalVfs; ivf++ {
			isAllocated, err := isVfAllocated(pf, uint(ivf))
			if err != nil {
				log.Printf("Error checking allocation: %s\n", err)
				continue
			}

			var vf VF
			if err := setVFConfig(pf, uint(ivf), &vf); err != nil {
				log.Printf("ERROR: Failed to add vf[%d] for pf[%s]: %s\n",
					ivf, pf, err)
				continue
			}
			vf.Allocated = isAllocated
			tmpNodePf.UsedTxRate += vf.MinTxRate
			tmpNodePf.CapacityVFs++
			if vf.Allocated {
				tmpNodePf.UsedVFs++
			}
			tmpNodePf.VFs = append(tmpNodePf.VFs, &vf)
		}
		if len(userConfig.PfMaxBandwidth) != len(userConfig.PfNetdevices) {
			log.Println("WARNING: not setting capacityTxRate b/c len(pfNetdevices) != len(pfMaxBandwidth)")
		} else {
			tmpNodePf.CapacityTxRate = userConfig.PfMaxBandwidth[0]
		}
		nodePfs = append(nodePfs, &tmpNodePf)
	}
	return nodePfs, nil
}

func CreateServer(port string, userConfig UserConfig) *http.Server {
	http.HandleFunc("/getpfs", func(w http.ResponseWriter, r *http.Request) {
		pfs, err := getNodeData(userConfig)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s", err)
			return
		}
		data, err := json.Marshal(pfs)
		if err != nil {
			fmt.Fprintf(w, "ERROR: failed to marshal pfs: %s", err)
			return
		}

		fmt.Fprintf(w, "%s", data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, message, version)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: nil}

	log.Println("Starting Server on port: " + port)
	return server
}
