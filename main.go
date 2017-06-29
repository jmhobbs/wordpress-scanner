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
)

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
	r.HandleFunc("/plugin/{plugin}/{version}", GetPlugin).Methods("GET")
	http.Handle("/", r)

	http.ListenAndServe("127.0.0.1:9090", r)
}

func GetPlugin(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	found := false

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vars["plugin"]))
		if b == nil {
			return nil
		}

		s := b.Get([]byte(vars["version"]))
		if s == nil {
			return nil
		}

		found = true
		w.Write(s)

		return nil
	})

	if found {
		return
	}

	scan, err := scanPlugin(vars["plugin"], vars["version"])
	if err != nil {
		panic(err)
	}

	//json.NewEncoder(w).Encode(scan)
	s, err := json.Marshal(scan)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(vars["plugin"]))
		if err != nil {
			log.Printf("Failed to create bucket: %s\n", err)
			return err
		}

		err = b.Put([]byte(vars["version"]), s)
		if err != nil {
			log.Printf("Failed to write scan: %s\n", err)
		}

		return nil
	})

	w.Write(s)
}

func scanPlugin(plugin, version string) (*Scan, error) {
	scan := NewScan(plugin, version)

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
			scan.addErrored(f.Name, err)
			continue
		}

		hash, err := GetHash(r)
		if err != nil {
			scan.addErrored(f.Name, err)
			continue
		}
		r.Close()

		scan.addHashed(f.Name, hash)
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
