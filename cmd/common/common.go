package common

import (
	"flag"
	"fmt"
	"os"

	"producer-consumer/pkg/build"
)

func ShowVersion() {
	versionFlag := flag.Bool("version", false, "prints the build version")
	flag.Parse()

	// If the version flag is passed, print the version and exit
	if *versionFlag {
		fmt.Printf("Version: %s\n", build.Version)
		os.Exit(0)
	}
}
