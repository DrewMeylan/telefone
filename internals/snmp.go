package internals

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/gosnmp/gosnmp"
)

// oidScanner holds configuration for scanning a subnet via SNMP.
type oidScanner struct {
	Subnet      *net.IPNet
	Version     string
	Community   string
	OID         string
	Concurrency int
	Timeout     time.Duration
}

// Scan performs the subnet scan, checking each IP for the presence of the given OID.
func (s *oidScanner) Scan() ([]net.IP, error) {
	var responsiveIPs []net.IP
	var wg sync.WaitGroup
	var mu sync.Mutex

	ips := helpers.enumerateIPs(s.Subnet)

	// Step 1 & 2: Ping sweep with port check
	pingJobs := make(chan net.IP, len(ips))
	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range pingJobs {
				if helpers.isHostUp(ip, s.Timeout) && helpers.isPortOpen(ip, 161, s.Timeout) {
					mu.Lock()
					responsiveIPs = append(responsiveIPs, ip)
					mu.Unlock()
				}
			}
		}()
	}
	for _, ip := range ips {
		pingJobs <- ip
	}
	close(pingJobs)
	wg.Wait()

	// Step 3 - 5: SNMP check for OID presence
	var matchedIPs []net.IP
	jobs := make(chan net.IP, len(responsiveIPs))
	wg = sync.WaitGroup{}
	mu = sync.Mutex{}

	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				if s.checkOID(ip) {
					mu.Lock()
					matchedIPs = append(matchedIPs, ip)
					mu.Unlock()
				}
			}
		}()
	}
	for _, ip := range responsiveIPs {
		jobs <- ip
	}
	close(jobs)
	wg.Wait()

	return matchedIPs, nil
}

// checkOID queries the OID on the given IP using the provided SNMP version and community string.
func (s *oidScanner) checkOID(ip net.IP) bool {
	params := &gosnmp.GoSNMP{
		Target:    ip.String(),
		Port:      161,
		Community: s.Community,
		Version:   getSNMPVersion(s.Version),
		Timeout:   s.Timeout,
		Retries:   1,
	}
	err := params.Connect()
	if err != nil {
		return false
	}
	defer params.Conn.Close()

	result, err := params.Get([]string{s.OID})
	if err != nil || len(result.Variables) == 0 {
		return false
	}
	return true
}

// getSNMPVersion returns the gosnmp SNMP version enum from string.
func getSNMPVersion(ver string) gosnmp.SnmpVersion {
	switch strings.ToLower(ver) {
	case "1":
		return gosnmp.Version1
	case "3":
		// SNMPv3 not implemented in this version
		log.Fatal("SNMPv3 not supported in this implementation")
		return gosnmp.Version3
	default:
		return gosnmp.Version2c
	}
}

func main() {
	var subnetStr, version, community, oid string
	var concurrency int
	var timeout time.Duration

	flag.StringVar(&subnetStr, "subnet", "", "IP subnet in CIDR notation (e.g., 192.168.1.0/24)")
	flag.StringVar(&version, "version", "2c", "SNMP version: 1, 2c, or 3")
	flag.StringVar(&community, "community", "public", "SNMPv2 community string")
	flag.StringVar(&oid, "oid", "", "OID to search for")
	flag.IntVar(&concurrency, "concurrency", 50, "Number of concurrent workers")
	flag.DurationVar(&timeout, "timeout", 2*time.Second, "Timeout for ping/SNMP operations")
	flag.Parse()

	if subnetStr == "" || oid == "" {
		fmt.Println("Usage: go run main.go -subnet=<cidr> -version=2c -community=public -oid=<oid>")
		os.Exit(1)
	}

	subnet, err := helpers.parseCIDR(subnetStr)
	if err != nil {
		log.Fatalf("Invalid subnet: %v", err)
	}

	scanner := &oidScanner{
		Subnet:      subnet,
		Version:     version,
		Community:   community,
		OID:         oid,
		Concurrency: concurrency,
		Timeout:     timeout,
	}

	results, err := scanner.Scan()
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	fmt.Println("Matching IPs:")
	for _, ip := range results {
		fmt.Println(ip)
	}
}

