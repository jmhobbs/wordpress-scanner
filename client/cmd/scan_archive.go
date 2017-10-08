package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
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
			fmt.Println("You must provide the plugin name, the plugin version, and the archive to scan")
			os.Exit(1)
		} else if len(args) > 3 {
			fmt.Println("You gave too many arguments")
			os.Exit(1)
		}

		plugin := args[0]
		version := args[1]
		archive := args[2]

		scan, err := scanPlugin(plugin, version, archive)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(os.Stdout).Encode(scan)
	},
}

func scanPlugin(plugin, version string, archive string) (*shared.Scan, error) {
	scan := shared.NewScan(plugin, version)

	file, err := os.Open(archive)
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
	}()

	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	r, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if f.Name[len(f.Name)-1] == '/' && f.UncompressedSize64 == 0 {
			continue
		}

		r, err := f.Open()
		if err != nil {
			scan.AddErrored(f.Name, err)
			continue
		}

		hash, err := shared.GetHash(r)
		if err != nil {
			scan.AddErrored(f.Name, err)
			continue
		}
		r.Close()

		scan.AddHashed(f.Name, hash)
	}

	return scan, nil
}
