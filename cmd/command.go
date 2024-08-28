/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dnstapir/tapir"
	"github.com/miekg/dns"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var BumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Instruct TEM to bump the SOA serial of the RPZ zone",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "bump",
			Zone:    dns.Fqdn(tapir.GlobalCF.Zone),
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

var TemCmd = &cobra.Command{
	Use:   "tem",
	Short: "Prefix command to TEM, only usable via sub-commands",
}

var TemMqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "Prefix command to TEM MQTT, only usable via sub-commands",
}

var TemStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of TEM",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "status",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)

		fmt.Printf("TemStatus: %v\n", resp.TemStatus)
		if len(resp.TemStatus.ComponentStatus) != 0 {
			ts := resp.TemStatus
			var out = []string{"Component|Status|Last event|Counters|Error msg|NumFailures|LastFailure"}
			for k, v := range ts.ComponentStatus {
				out = append(out, fmt.Sprintf("%s|%s|%s|%d|%s|%d|%s", k, v, ts.TimeStamps[k].Format(tapir.TimeLayout), ts.Counters[k], ts.ErrorMsgs[k], ts.NumFailures, ts.LastFailure.Format(tapir.TimeLayout)))
			}
			fmt.Printf("%s\n", columnize.SimpleFormat(out))
		}
	},
}

var TemStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Instruct TEM to stop",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "stop",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

var TemMqttStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Instruct TEM MQTT Engine to start",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "mqtt-start",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

var TemMqttStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Instruct TEM MQTT Engine to stop",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "mqtt-stop",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

var TemMqttRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Instruct TEM MQTT Engine to restart",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "mqtt-restart",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

func init() {
	rootCmd.AddCommand(BumpCmd, TemCmd)
	TemCmd.AddCommand(TemStatusCmd, TemStopCmd, TemMqttCmd)
	TemMqttCmd.AddCommand(TemMqttStartCmd, TemMqttStopCmd, TemMqttRestartCmd)

	BumpCmd.Flags().StringVarP(&tapir.GlobalCF.Zone, "zone", "z", "", "Zone name")
}

func SendCommandCmd(data tapir.CommandPost) tapir.CommandResponse {
	_, buf, _ := api.RequestNG(http.MethodPost, "/command", data, true)

	var cr tapir.CommandResponse

	err := json.Unmarshal(buf, &cr)
	if err != nil {
		log.Fatalf("Error from json.Unmarshal: %v\n", err)
	}
	return cr
}
