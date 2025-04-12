package snmp

import (
	"fmt"
	"time"
	"net"
	"sync"

	"github.com/gosnmp/gosnmp"
)

// should Scan() methods produce JUST ip addresses or full target path

type snmpScanner struct {
	// This is essentially just a port scanner; search subnet for devices listening on 161
	Subnet  					*net.IPNet
	Port 							uint16
	Timeout 					time.Duration
	ConcurrentScans		int
	Results           []string	
}

func (s *snmpScanner) Scan(net.IPNet) ([]net.IP, error) {
	
	// We should add logic in to pingsweep the IP Subnet first so we can perform the port scan only on active hosts.
	// This will save a little time on the port scanner itself but substantial time on scanners with retries, timeouts, etc.

	var wg snc.WaitGroup // Review Waitgroup
	ipChan := make(chan string, s.ConcurrentScans) // Review
	s.Results = make(map[string]map[string]string) //Review
	mutex := sync.Mutex{} //Review

	for i := 0; i < s.ConcurrentScans; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				results := s.queryHost(ip) // This is a METHOD; we need to define this method within the context of snmpScanner
				if results != nil {
					mutex.Lock()
					s.Results[ip] = results
					mutex.Unlock()
				}
			}
		}()
	}
	for ip := range hosts(s.Subnet) {
		ipChan <- ip
	}
	close(ipChan)
	wg.Wait()
	return nil // This almost certainly is not correct. Needs to return []net.IP, err
}

type oidScanner struct {
	// Search subnet for devices containing a specific set of OIDs
	Subnet          *net.IPNet   // The subnet to scan (e.g., 192.168.1.0/24)
	Community       string       // SNMP community string (e.g., "public")
	Version         string       // SNMP version (e.g., "2c")
	OIDs            []string     // List of OIDs to scan for
	Timeout         time.Duration // Timeout per request
	Retries         int           // Number of retries for each request
	Port            uint16        // SNMP port (default 161)
	ConcurrentScans int           // Number of concurrent workers for scanning
	Results         map[string]map[string]string // Results[IP][OID] = value
	Verbose         bool          // Enable detailed logging/debugging
}

func (s *oidScanner) Scan(net.IPNet) ([]net.IP, error) {
	// MODIFY return more than just an error
	var wg sync.WaitGroup
	ipChan := make(chan string, s.ConcurrentScans)
	s.Results = make(map[string]map[string]string)
	mutex := sync.Mutex{}

	// Start worker goroutines
	for i := 0; i < s.ConcurrentScans; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				results := s.queryHost(ip)
				if results != nil {
					mutex.Lock()
					s.Results[ip] = results
					mutex.Unlock()
				}
			}
		}()
	}

	// Enqueue IPs
	for ip := range hosts(s.Subnet) {
		ipChan <- ip
	}
	close(ipChan)
	wg.Wait()
	return nil // This does not appear to be correct? 
}

func (s *oidScanner) queryHost(ip string) map[string]string {
	params := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      s.Port,
		Community: s.Community,
		Version:   gosnmp.Version2c,
		Timeout:   s.Timeout,
		Retries:   s.Retries,
	}

	if err := params.Connect(); err != nil {
		if s.Verbose {
			fmt.Printf("[!] Failed to connect to %s: %v\n", ip, err)
		}
		return nil
	}
	defer params.Conn.Close()

	result, err := params.Get(s.OIDs)
	if err != nil {
		if s.Verbose {
			fmt.Printf("[!] SNMP GET failed for %s: %v\n", ip, err)
		}
		return nil
	}

	response := make(map[string]string)
	for _, variable := range result.Variables {
		oid := variable.Name
		value := fmt.Sprintf("%v", variable.Value)
		response[oid] = value
		if s.Verbose {
			fmt.Printf("[+] %s - %s: %s\n", ip, oid, value)
		}
	}
	return response
}

