/*


Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net"
	"github.com/drewmeylan/telefone/internals"
	"github.com/spf13/cobra"
)

var (
	network  string
	oid 		 string
	community string
	version string
)

// scanOidCmd represents the scanOid command
var scanOidCmd = &cobra.Command{
	Use:   "scanOid",
	Short: "Scan devices in a subnet for a specific OID",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("network: ", network)	
		fmt.Println("oid", oid)
		fmt.Println("community", community)
		fmt.Println("version", version)

		ip_range := helpers.EnumerateHosts(net.IPNet(network))
		fmt.Println(ip_range)
	},
}

func init() {
	rootCmd.AddCommand(scanOidCmd)

	scanOidCmd.Flags().StringVarP(&version, "version", "v", "", "SNMP version")
	scanOidCmd.Flags().StringVarP(&oid, "oid", "o", "", "OID to search for")
	scanOidCmd.Flags().StringVarP(&community, "community", "c", "", "Community string")
	scanOidCmd.Flags().StringVarP(&network, "network", "n", "", "Network to scan (CIDR)")
		
	scanOidCmd.MarkFlagRequired("version")
	scanOidCmd.MarkFlagRequired("oid")
	scanOidCmd.MarkFlagRequired("community")
	scanOidCmd.MarkFlagRequired("network")
}
