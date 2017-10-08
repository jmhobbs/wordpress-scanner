package cmd

import (
	"fmt"
	"os"
)

func errorAndExit(message string) {
	fmt.Println("You must provide the plugin name, the plugin version, and the archive to scan")
	os.Exit(1)
}
