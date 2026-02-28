package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Params struct {
	Target  string
	Source  string
	Times   int
	Timeout time.Duration
}

type Result struct {
	Target  string
	Source  string
	Times   int
	Replies int
	RTTs    []time.Duration
	Error   string
}

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
	buffer := make([]byte, 1500)

	for seq := 1; seq <= p.Times; seq++ {

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  seq,
				Data: []byte("PING"),
			},
		}

		data, err := msg.Marshal(nil)
		if err != nil {
			continue
		}

		start := time.Now()

		_, err = conn.WriteTo(data, ipAddr)
		if err != nil {
			continue
		}

		_ = conn.SetReadDeadline(time.Now().Add(p.Timeout))

		for {
			n, peer, err := conn.ReadFrom(buffer)
			if err != nil {
				break // timeout
			}

			// pastikan dari target yang sama
			if peer.String() != ipAddr.String() {
				continue
			}

			rm, err := icmp.ParseMessage(1, buffer[:n])
			if err != nil {
				continue
			}

			if rm.Type != ipv4.ICMPTypeEchoReply {
				continue
			}

			body, ok := rm.Body.(*icmp.Echo)
			if !ok {
				continue
			}

			// VALIDASI ID & SEQ
			if body.ID == id && body.Seq == seq {
				rtt := time.Since(start)
				result.Replies++
				result.RTTs = append(result.RTTs, rtt)

				fmt.Printf("Reply from %s seq=%d time=%v\n",
					ipAddr.String(), seq, rtt)

				break
			}
		}

		time.Sleep(500 * time.Millisecond) // interval seperti ping
	}

	return result
}