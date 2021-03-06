package main

import (
	"fmt"
	"os"
)

// BuildDate is the date & time when the binary was built
var BuildDate string

// GitCommit is the commit hash when the binary was built
var GitCommit string

// Version is the version of the binary
var Version string

// PrintVersionAndExit prints the version and exits
func printVersionAndExit() {
	fmt.Printf("Version: %s - Commit: %s - Date: %s\n",
		Version, GitCommit, BuildDate)
	os.Exit(0)
}
