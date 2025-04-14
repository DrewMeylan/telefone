/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// scanOidCmd represents the scanOid command
var scanOidCmd = &cobra.Command{
	Use:   "scanOid",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scanOid called")

	},
}

func init() {
	rootCmd.AddCommand(scanOidCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanOidCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	
	scanOidCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanOidCmd.Flags().StringVarP("version", "v",  "version") 
	scanOidCmd.Flags().StringVarP("oid", "o", "", "OID to search for") 
	scanOidCmd.Flags().StringVarP("comunity", "c", "", "Community string")
	scanOidCmd.Flags().StringVarP("network", "n", "", "Network to scan (CIDR)")
	
	scanOidCmd.MarkFlagsRequired("version")
	scanOidCmd.MarkFlagsRequired("oid")
	scanOidCmd.MarkFlagsRequired("community")
	scanOidCmd.MarkFlagsRequired("network")
}
