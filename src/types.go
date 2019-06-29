package src

type PF struct {
	UsedTxRate     uint  `json:"used_tx_rate"`
	CapacityTxRate uint  `json:"capacity_tx_rate"`
	UsedVFs        uint  `json:"used_vfs"`
	CapacityVFs    uint  `json:"capacity_vfs"`
	VFs            []*VF `json:"vfs"`
}

type VF struct {
	VFNumber   uint   `json:"vf"`
	MAC        string `json:"mac"`
	VLAN       uint   `json:"vlan"`
	QoS        uint   `json:"qos"`
	VLanProto  string `json:"vlan_proto"`
	SpoofCheck string `json:"spoof_check"`
	Trust      string `json:"trust"`
	LinkState  string `json:"link_state"`
	MinTxRate  uint   `json:"min_tx_rate"`
	MaxTxRate  uint   `json:"max_tx_rate"`
	VGTPlus    string `json:"vgt_plus"`
	RateGroup  uint   `json:"rate_group"`
	Allocated  bool   `json:"allocated"`
}

type UserConfig struct {
	Mode           string   `json:"mode"`
	PfNetdevices   []string `json:"pfNetdevices"`
	PfMaxBandwidth []uint   `json:"pfMaxBandwidth"`
}
