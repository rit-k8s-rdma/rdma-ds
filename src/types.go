package src

import "fmt"

type PfNetDevice struct {
	Name           string
	MaxSendingRate uint
}

type SystemConfig struct {
	PfNetDevices []PfNetDevice
}

func (s *SystemConfig) GetDeviceNames() []string {
	deviceNames := make([]string, len(s.PfNetDevices), len(s.PfNetDevices))
	for ipf, pf := range s.PfNetDevices {
		deviceNames[ipf] = pf.Name
	}
	return deviceNames
}

func (s *SystemConfig) GetDeviceSendingRate(device string) (uint, error) {
	for _, pf := range s.PfNetDevices {
		if pf.Name == device {
			return pf.MaxSendingRate, nil
		}
	}
	return 0, fmt.Errorf("Device not found")
}
