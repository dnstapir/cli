/*
 * Johan Stenstam, johani@johani.org
 */
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/dnstapir/tapir"
)

var newapi bool

const timelayout = "2006-01-02 15:04:05"

var ServerName string = "PLACEHOLDER"

var TemPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send an API ping request and present the response",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			log.Fatal("ping must have no arguments")
		}

		pr, err := tapir.GlobalCF.Api.SendPing(tapir.GlobalCF.PingCount, false)
		if err != nil {
			log.Fatalf("Error from SendPing: %v", err)
		}

		uptime := time.Now().Sub(pr.BootTime).Round(time.Second)
		if tapir.GlobalCF.Verbose {
			fmt.Printf("%s from %s @ %s (version %s): pings: %d, pongs: %d, uptime: %v time: %s, client: %s\n",
				pr.Msg, pr.Daemon, pr.ServerHost, pr.Version, pr.Pings,
				pr.Pongs, uptime, pr.Time.Format(timelayout), pr.Client)
		} else {
			fmt.Printf("%s: pings: %d, pongs: %d, uptime: %v, time: %s\n",
				pr.Msg, pr.Pings, pr.Pongs, uptime, pr.Time.Format(timelayout))
		}
	},
}

var DaemonApiCmd = &cobra.Command{
	Use:   "api",
	Short: "request a TEM api summary",
	Long:  `Query TEM for the provided API endpoints and print that out in a (hopefully) comprehensible fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			log.Fatal("api must have no arguments")
		}
		tapir.GlobalCF.Api.ShowApi()
	},
}

func init() {
	rootCmd.AddCommand(DaemonApiCmd)
	TemCmd.AddCommand(TemPingCmd)

	TemPingCmd.Flags().IntVarP(&tapir.GlobalCF.PingCount, "count", "c", 0, "#pings to send")
	TemPingCmd.Flags().BoolVarP(&newapi, "newapi", "n", false, "use new api client")
}
