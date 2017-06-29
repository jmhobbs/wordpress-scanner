package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmhobbs/wordpress-scanner/shared"
)

func main() {
	plugin := os.Args[1]
	version := os.Args[2]
	path := os.Args[3]

	scan := shared.NewScan(plugin, version)

	err := filepath.Walk(path, scanFile(scan))
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(scan)
}

func scanFile(scan *shared.Scan) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".php") {
			f, err := os.Open(path)
			if err != nil {
				scan.AddErrored(path, err)
				return nil
			}

			hash, err := shared.GetHash(f)
			if err != nil {
				scan.AddErrored(path, err)
				return nil
			}

			scan.AddHashed(path, hash)
		}

		return nil
	}
}
