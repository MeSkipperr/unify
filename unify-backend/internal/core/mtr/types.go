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
	JSON   bool // default true (karena kita parse)

	IPv4Only bool // -4
	IPv6Only bool // -6
}

// Struktur output JSON MTR (disederhanakan)
type Result struct {
	Host      string      `json:"host"`
	TotalHops int         `json:"total_hops"`
	Reachable bool        `json:"reachable"`
	MaxLoss   float64     `json:"max_loss"`
	AvgRTT    float64     `json:"avg_rtt"`
	Hops      []HopResult `json:"hops"`
}

type HopResult struct {
	Hop    int     `json:"hop"`
	IP     string  `json:"ip"`
	DNS   string  `json:"dns"`
	Loss   float64 `json:"loss"`
	AvgRTT float64 `json:"avg_rtt"`
}

type MtrJSON struct {
	Report struct {
		MTR  struct {
			Dst   string `json:"dst"`
			Src   string `json:"src"`
			Pings int    `json:"pings"`
		} `json:"mtr"`
		Hubs []struct {
			Count int     `json:"count"`
			Host  string  `json:"host"`
			Loss  float64 `json:"Loss%"`
			Avg   float64 `json:"Avg"`
		} `json:"hubs"`
	} `json:"report"`
}
