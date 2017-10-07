package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/jmhobbs/wordpress-scanner/shared"
)

type Version struct {
	Version string `json:"version"`
	Files   []shared.File `json:"files"`
}

type VersionList struct {
	Plugin string `json:"plugin"`
	Versions []Version `json:"versions"`
}

var db *bolt.DB

func main() {
	var err error
	log.Println("Opening database")
	db, err = bolt.Open("plugins.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/plugin/{plugin}/{version}/diff", DiffPlugin).Methods("POST")
	r.HandleFunc("/plugin/{plugin}/{version}", GetPlugin).Methods("GET")
	r.HandleFunc("/plugin/{plugin}", ListPluginVersions).Methods("GET")
	r.HandleFunc("/plugin", ListPlugins).Methods("GET")
	http.Handle("/", r)

	http.ListenAndServe("127.0.0.1:9090", r)
}

func DiffPlugin(w http.ResponseWriter, req *http.Request) {
	vars    := mux.Vars(req)
	plugin  := vars["plugin"]
	version := vars["version"]

	var scan shared.Scan
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&scan)

	if err != nil {
		panic(err)
	}

	defer req.Body.Close()

	referenceScan, err := lookupOrScanPlugin(plugin, version)

	if err != nil {
		panic(err)
	}

	diff := referenceScan.Diff(&scan)
	d, err := json.Marshal(diff)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func GetPlugin(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	plugin := vars["plugin"]
	version := vars["version"]

	scan, err := lookupOrScanPlugin(plugin, version)

	if err != nil {
		panic(err)
	}

	s, err := json.Marshal(scan)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(s)
}

func ListPlugins(w http.ResponseWriter, req *http.Request) {
	plugins := make([]string, 0)

	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucket []byte,_ *bolt.Bucket) error {
			plugins = append(plugins, string(bucket))
			return nil
		})
	})

	output, err := json.Marshal(plugins)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func ListPluginVersions(w http.ResponseWriter, req *http.Request) {
	plugin      := mux.Vars(req)["plugin"]
	versionList := VersionList{Plugin: plugin, Versions: make([]Version, 0)}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(plugin))

		if b == nil {
			return nil
		}

		b.ForEach(func(version, s []byte) error {
			var scan shared.Scan
			err := json.Unmarshal(s, &scan)

			if err != nil {
				log.Printf("error: %s", err)
				return nil
			}

			versionList.Versions = append(versionList.Versions, Version{Version: string(version), Files: scan.Files})

			return nil
		})

		return nil
	})

	output, err := json.Marshal(versionList)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func scanPlugin(plugin, version string) (*shared.Scan, error) {
	scan := shared.NewScan(plugin, version)

	log.Println("Downloading Plugin Zip to TempFile")
	tmp, b, err := downloadPluginFile(plugin, version)
	if err != nil {
		return nil, err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	log.Printf("Downloaded to %s\n", tmp.Name())

	log.Println("Opening zip file")
	r, err := zip.NewReader(tmp, b)
	if err != nil {
		return nil, err
	}

	log.Println("Hashing files")
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

// Returns a zip file in a temporary location.
// Caller is responsible for closing and removing.
// i.e. os.Remove(tmpfile.Name())
func downloadPluginFile(name, version string) (*os.File, int64, error) {
	filename := strings.Join([]string{name, version, "zip"}, ".")

	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		return nil, 0, err
	}

	resp, err := http.Get("https://downloads.wordpress.org/plugin/" + filename)
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, 0, err
	}
	defer resp.Body.Close()

	written, err := io.Copy(tmpfile, resp.Body)
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, 0, err
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, 0, err
	}

	return tmpfile, written, nil
}

func lookupOrScanPlugin(plugin string, version string) (*shared.Scan, error) {
	var s []byte
	found := false

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(plugin))
		if b == nil {
			return nil
		}

		s = b.Get([]byte(version))
		if s == nil {
			return nil
		}

		found = true

		return nil
	})

	if found {
		var scan shared.Scan
		err := json.Unmarshal(s, &scan)

		if err != nil {
			panic(err)
		}

		return &scan, nil
	}

	scan, err := scanPlugin(plugin, version)
	if err != nil {
		panic(err)
	}

	s, err = json.Marshal(scan)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(plugin))
		if err != nil {
			log.Printf("Failed to create bucket: %s\n", err)
			return err
		}

		err = b.Put([]byte(version), s)
		if err != nil {
			log.Printf("Failed to write scan: %s\n", err)
		}

		return nil
	})

	log.Println("Returning from the bottom")
	return scan, nil
}
