package main

import (
	"io"

	"github.com/spaolacci/murmur3"
)

type File struct {
	Path  string `json:'path'`
	Error error
	Hash  uint32 `json:'hash'`
}

type Scan struct {
	Plugin  string `json:'plugin'`
	Version string `json:'version'`
	Files   []File `json:'files'`
}

func NewScan(plugin, version string) *Scan {
	return &Scan{plugin, version, []File{}}
}

func (s *Scan) addHashed(path string, hash uint32) {
	s.Files = append(s.Files, File{path, nil, hash})
}

func (s *Scan) addErrored(path string, err error) {
	s.Files = append(s.Files, File{path, err, 0})
}

func GetHash(r io.Reader) (uint32, error) {
	chunk := make([]byte, 100) // TODO: Tune this.
	hash := murmur3.New32()
	for {
		n, err := r.Read(chunk)
		if err != nil {
			if err == io.EOF {
				if n == 0 {
					break
				}
			} else {
				return 0, err
			}
		}
		hash.Write(chunk[:n])
	}
	return hash.Sum32(), nil
}
