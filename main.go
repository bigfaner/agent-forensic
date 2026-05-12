package main

import "github.com/user/agent-forensic/cmd"

// Version is set via -ldflags at build time. Defaults to "dev".
var Version = "dev"

func main() {
	cmd.SetVersion(Version)
	cmd.Execute()
}
