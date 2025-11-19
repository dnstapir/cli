/*
 * Copyright 2024 Johan Stenstam, johan.stenstam@internetstiftelsen.se
 */

package main

import (
    "fmt"
	"dnstapir-cli/cmd"
)

var name    = "BAD-BUILD"
var version = "BAD-BUILD"
var commit  = "BAD-BUILD"

func main() {
    fmt.Printf("%s, version %s, commit %s\n", name, version, commit)
	cmd.Execute()
}
