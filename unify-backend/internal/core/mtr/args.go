package mtr

import "strconv"

func buildArgs(cfg Config) []string {
	args := []string{}

	// Output
	if cfg.JSON {
		args = append(args, "--json")
	}

	args = append(args, "-c", strconv.Itoa(cfg.Count))

	// DNS
	if !cfg.UseDNS {
		args = append(args, "--no-dns")
	}

	// IP version
	if cfg.IPv4Only {
		args = append(args, "-4")
	}
	if cfg.IPv6Only {
		args = append(args, "-6")
	}

	// Source
	if cfg.SourceIP != "" {
		args = append(args, "-a", cfg.SourceIP)
	}
	if cfg.Interface != "" {
		args = append(args, "-I", cfg.Interface)
	}

	// TTL
	if cfg.FirstTTL > 0 {
		args = append(args, "-f", strconv.Itoa(cfg.FirstTTL))
	}
	if cfg.MaxTTL > 0 {
		args = append(args, "-m", strconv.Itoa(cfg.MaxTTL))
	}

	// Timing
	if cfg.Interval > 0 {
		args = append(args, "-i", strconv.Itoa(cfg.Interval))
	}
	if cfg.Timeout > 0 {
		args = append(args, "-Z", strconv.Itoa(cfg.Timeout))
	}

	// Packet
	if cfg.PacketSize > 0 {
		args = append(args, "-s", strconv.Itoa(cfg.PacketSize))
	}

	// Protocol
	switch cfg.Protocol {
	case ProtocolTCP:
		args = append(args, "--tcp")
	case ProtocolUDP:
		args = append(args, "--udp")
	default:
		// ICMP default → no flag
	}

	// Port
	if cfg.Port != nil {
		args = append(args, "-P", strconv.Itoa(*cfg.Port))
	}
	if cfg.LocalPort != nil {
		args = append(args, "-L", strconv.Itoa(*cfg.LocalPort))
	}

	// Destination
	args = append(args, cfg.DestHost)

	return args
}
