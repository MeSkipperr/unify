package ip

import (
	"fmt"
	"net"
	"runtime"
)

type LocalIPInfo struct {
	Exists       bool
	IPAddress    string
	Interface    string
	Netmask      string
	CIDR         string
	IPVersion    int // 4 or 6
	IsLoopback   bool
	IsUp         bool
}

// HasLocalIP mengecek apakah IP ada di interface lokal
func CheckLocalIp(ip string) (*LocalIPInfo, error) {
	switch runtime.GOOS {
	case "linux":
		return getLocalIPInfoLinux(ip)
	case "windows":
		return getLocalIPInfoWindows(ip)
	default:
		return nil, fmt.Errorf("unsupported OS")
	}
}


func getLocalIPInfoLinux(ip string) (*LocalIPInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.String() == ip {
				ones, _ := ipNet.Mask.Size()

				return &LocalIPInfo{
					Exists:     true,
					IPAddress:  ip,
					Interface:  iface.Name,
					Netmask:   net.IP(ipNet.Mask).String(),
					CIDR:      fmt.Sprintf("%s/%d", ip, ones),
					IPVersion: func() int {
						if ipNet.IP.To4() != nil {
							return 4
						}
						return 6
					}(),
					IsLoopback: iface.Flags&net.FlagLoopback != 0,
					IsUp:       iface.Flags&net.FlagUp != 0,
				}, nil
			}

		}
	}

	return &LocalIPInfo{Exists: false}, nil
}


func getLocalIPInfoWindows(ip string) (*LocalIPInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.String() == ip {
				ones, _ := ipNet.Mask.Size()

				return &LocalIPInfo{
					Exists:     true,
					IPAddress:  ip,
					Interface:  iface.Name,
					Netmask:   net.IP(ipNet.Mask).String(),
					CIDR:      fmt.Sprintf("%s/%d", ip, ones),
					IPVersion: func() int {
						if ipNet.IP.To4() != nil {
							return 4
						}
						return 6
					}(),
					IsLoopback: iface.Flags&net.FlagLoopback != 0,
					IsUp:       iface.Flags&net.FlagUp != 0,
				}, nil
			}
		}
	}

	return &LocalIPInfo{Exists: false}, nil
}
