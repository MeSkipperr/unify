package mtr

type Protocol string

const (
	ProtocolICMP Protocol = "icmp"
	ProtocolTCP  Protocol = "tcp"
	ProtocolUDP  Protocol = "udp"
)

type Config struct {
	// Target
	DestHost string // hostname atau IP (required)

	// Network
	SourceIP  string // -a, --address (optional)
	Interface string // -I, --interface (optional)

	// Protocol
	Protocol  Protocol // default: icmp
	Port      *int     // -P (tcp/udp only)
	LocalPort *int     // -L (udp source port)

	// Probe behavior
	Count    int // -c, default 10
	Interval int // -i, seconds, default mtr
	Timeout  int // -Z, seconds
	FirstTTL int // -f
	MaxTTL   int // -m

	// Packet
	PacketSize int // -s

	// Output
	UseDNS bool // default true
	JSON   bool // default true

	IPv4Only bool // -4
	IPv6Only bool // -6
}

type MtrType struct {
	Src        string `json:"src"`        // Source hostname or IP where the test starts
	Dst        string `json:"dst"`        // Destination IP or hostname being tested
	Tos        int    `json:"tos"`        // Type of Service (DSCP/TOS value), usually 0
	Tests      int    `json:"tests"`      // Number of packets sent per hop
	Psize      string `json:"psize"`      // Packet payload size in bytes
	BitPattern string `json:"bitpattern"` // Payload bit pattern used in the ICMP packets
	TotalHops  int
	Reachable  bool //default false
	AvgRTT     float64
}

type HopsType struct {
	Count int    `json:"count"` // Hop number in the network path
	Host  string `json:"host"`  // Router IP or hostname at this hop
	Dns   string
	Loss  float64 `json:"Loss%"` // Packet loss percentage at this hop
	Snt   int     `json:"Snt"`   // Total packets sent to this hop
	Last  float64 `json:"Last"`  // Round-trip time of the last packet (ms)
	Avg   float64 `json:"Avg"`   // Average round-trip time (ms)
	Best  float64 `json:"Best"`  // Fastest round-trip time (ms)
	Worst float64 `json:"Wrst"`  // Slowest round-trip time (ms)
	StDev float64 `json:"Stdev"` // RTT standard deviation (jitter indicator)
}

type MtrResultJson struct {
	Report struct {
		Result    MtrType    `json:"mtr"`  // MTR metadata (source, destination, tests, etc.)
		HopResult []HopsType `json:"hubs"` // List of hop results (each router hop)
	} `json:"report"`
}
