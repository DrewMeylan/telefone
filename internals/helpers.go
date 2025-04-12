package helpers

import (
	"net"
)

// hosts generates all usable IPs in the subnet.
func hosts(network *net.IPNet) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for ip := network.IP.Mask(network.Mask); network.Contains(ip); incIP(ip) {
			ipCopy := make(net.IP, len(ip))
			copy(ipCopy, ip)
			// Skip network and broadcast addresses
			if !ipCopy.Equal(network.IP) && !ipCopy.Equal(lastIP(network)) {
				ch <- ipCopy.String()
			}
		}
	}()
	return ch
}

// incIP increments an IP address.
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

// lastIP returns the broadcast IP of the subnet.
func lastIP(n *net.IPNet) net.IP {
	ip := make(net.IP, len(n.IP))
	copy(ip, n.IP)
	for i := range ip {
		ip[i] |= ^n.Mask[i]
	}
	return ip
}
