package internals

import (
	"net"
	"time"
	"fmt"

	"github.com/go-ping/ping"
)

// enumerateIPs expands a CIDR subnet into a slice of IPs.
func enumerateIPs(subnet *net.IPNet) []net.IP {
	var ips []net.IP
	for ip := subnet.IP.Mask(subnet.Mask); subnet.Contains(ip); incIP(ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy)
	}
	// Remove network and broadcast addresses
	if len(ips) > 2 {
		return ips[1 : len(ips)-1]
	}
	return ips
}

// incIP increments an IP address by 1.
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// isHostUp checks if a host is reachable via ping.
func isHostUp(ip net.IP, timeout time.Duration) bool {
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		return false
	}
	pinger.Count = 1
	pinger.Timeout = timeout
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		return false
	}
	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

// isPortOpen checks if a specific port is open on the host.
func isPortOpen(ip net.IP, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", ip.String(), port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// parseCIDR parses and validates a CIDR string.
func parseCIDR(subnet string) (*net.IPNet, error) {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, err
	}
	return ipnet, nil
}

