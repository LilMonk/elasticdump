package main

import (
	"fmt"
	"os"

	"github.com/lilmonk/elasticdump/cmd"
)

// version is set during build time using ldflags
var version = "dev"

func main() {
	// Set version in cmd package
	cmd.SetVersion(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
