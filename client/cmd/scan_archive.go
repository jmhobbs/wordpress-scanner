package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/jmhobbs/wordpress-scanner/shared"
)

var scanArchiveCmd = &cobra.Command{
	Use:   "scan-archive <plugin_name> <plugin_version> <plugin_archive>",
	Short: "Scans a plugin archive for corruption",
	Long:  "",
	Run: func (cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			errorAndExit("You must provide the plugin name, the plugin version, and the archive to scan")
		} else if len(args) > 3 {
			errorAndExit("You gave too many arguments")
		}

		plugin := args[0]
		version := args[1]
		archive := args[2]

		file, err := os.Open(archive)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			file.Close()
		}()

		scan, err := shared.NewScanFromFile(plugin, version, file)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(os.Stdout).Encode(scan)
	},
}
