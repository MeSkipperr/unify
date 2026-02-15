package ping

import (
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Params struct {
	Target  string // destination IP / hostname
	Source  string // optional source IP (linux)
	Times   int    // number of echo requests
	Timeout time.Duration
}

type Result struct {
	Target  string          `json:"target"`
	Source  string          `json:"source,omitempty"`
	Times   int             `json:"times"`
	Replies int             `json:"replies"`
	RTTs    []time.Duration `json:"rtts"`
	Error   string          `json:"error,omitempty"`
}

// Ping sends ICMP Echo Requests using raw packets
func Ping(p Params) Result {
	if p.Times <= 0 {
		p.Times = 1
	}
	if p.Timeout <= 0 {
		p.Timeout = 2 * time.Second
	}

	result := Result{
		Target: p.Target,
		Source: p.Source,
		Times:  p.Times,
		RTTs:   make([]time.Duration, 0),
	}

	ipAddr, err := net.ResolveIPAddr("ip4", p.Target)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	conn, err := icmp.ListenPacket("ip4:icmp", p.Source)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer conn.Close()

	id := os.Getpid() & 0xffff

	for i := 1; i <= p.Times; i++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  i,
				Data: []byte("PING"),
			},
		}

		data, _ := msg.Marshal(nil)
		start := time.Now().UTC()

		_, err := conn.WriteTo(data, ipAddr)
		if err != nil {
			continue
		}

		_ = conn.SetReadDeadline(time.Now().UTC().Add(p.Timeout))
		buf := make([]byte, 1500)

		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buf[:n])
		if err != nil {
			continue
		}

		if rm.Type == ipv4.ICMPTypeEchoReply {
			result.Replies++
			result.RTTs = append(result.RTTs, time.Since(start))
		}
	}

	return result
}
