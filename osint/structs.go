package osint

// HostStruct holds temporary values before converting to JSON
type HostStruct struct {
	IPAddress     string
	Hostname      []string
	OSINTResponse OSINTResponses
}

// OSINTResponses holds responses from OSINT resources
type OSINTResponses struct {
	Shodan     OSINTInfo
	Binaryedge OSINTInfo
	Censys     OSINTInfo
	Zoomeye    OSINTInfo
	Onyphe     OSINTInfo
	Spyse      OSINTInfo
}

// OSINTInfo is a generic struct for OSINT information
type OSINTInfo struct {
	OpenPorts []int
}
