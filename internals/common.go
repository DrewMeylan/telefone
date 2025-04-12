package common

import (
	"fmt"
)

func discover [T Scanner](my_scanner Scanner, net *net.IPNet) ([]net.IP, error) {
	results, err := my_scanner.Scan(net)
	if err != nil {
		return nil, fmt.Errorf("Error encountered: %w", err)
	}

	return results, nil
}
