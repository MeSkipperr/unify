package utils

import (
	"errors"
	"net"
	"strings"
)

func NormalizeIPv4(ip string) (string, error) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return "", errors.New("invalid IPv4 address")
	}
	return parsed.String(), nil
}

func NormalizeMac(mac string) (string, error) {
	mac = strings.ToUpper(strings.TrimSpace(mac))
	_, err := net.ParseMAC(mac)
	if err != nil {
		return "", errors.New("invalid MAC address")
	}
	return mac, nil
}
