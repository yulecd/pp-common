package trace

import (
	"fmt"
	"log"
	"net"
)

// GetLocalIP
// return empty string if error
// NOTE: hex format ip
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("get local ip failed, error: %s\n", err)
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok {
			if isIPUseful(ipNet.IP) {
				bytes := ipNet.IP.To4()
				return fmt.Sprintf("%02x%02x%02x%02x", bytes[0], bytes[1], bytes[2], bytes[3])
			}
		}
	}

	log.Println("no valid ipv4")
	return ""
}

// GetLocalIPDotFormat
// return empty string if error
func GetLocalIPDotFormat() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("get local ip failed, error: %s\n", err)
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok {
			if isIPUseful(ipNet.IP) {
				return ipNet.IP.String()
			}
		}
	}

	log.Println("no valid ipv4")
	return ""
}

// private IPv4
// Class        Starting IPAddress    Ending IP Address    # of Hosts
// A            10.0.0.0              10.255.255.255       16,777,216
// B            172.16.0.0            172.31.255.255       1,048,576
// C            192.168.0.0           192.168.255.255      65,536
// Link-local-u 169.254.0.0           169.254.255.255      65,536
// Link-local-m 224.0.0.0             224.0.0.255          256
// Local        127.0.0.0             127.255.255.255      16777216
func isIPUseful(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return false
	}
	return true
}
