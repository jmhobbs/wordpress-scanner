package shared

import (
	"archive/zip"
	"os"
	"strings"
)

type Diff struct {
	Plugin  *DiffName  `json:"plugin,omitempty"`
	Version *DiffName  `json:"version,omitempty"`
	Files   []DiffFile `json:"files"`
}

type DiffFile struct {
	Path          string  `json:"path"`
	ReferenceHash *uint32 `json:"reference"`
	GivenHash     *uint32 `json:"given"`
}

type DiffName struct {
	ReferenceName *string `json:"reference"`
	GivenName     *string `json:"given"`
}

type File struct {
	Path  string `json:"path"`
	Error error  `json:"error,omitempty"`
	Hash  uint32 `json:"hash"`
}

type Scan struct {
	Plugin  string `json:"plugin"`
	Version string `json:"version"`
	Files   []File `json:"files"`
}

func NewScan(plugin, version string) *Scan {
	return &Scan{plugin, version, []File{}}
}

func NewScanFromFile(plugin, version string, file *os.File) (*Scan, error) {
	scan := NewScan(plugin, version)

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

		scan.Scan(f.Name)
	}

	return scan, nil
}

func (s *Scan) AddHashed(path string, hash uint32) {
	s.Files = append(s.Files, File{path, nil, hash})
}

func (s *Scan) AddErrored(path string, err error) {
	s.Files = append(s.Files, File{path, err, 0})
}

func (s *Scan) Diff(other *Scan) *Diff {
	diff := Diff{nil, nil, []DiffFile{}}
	diff.AddName(s.Plugin, other.Plugin)
	diff.AddVersion(s.Version, other.Version)

	for _, otherFile := range other.Files {
		found := false

		for _, file := range s.Files {
			if file.Path == otherFile.Path {
				if file.Hash != otherFile.Hash {
					diff.AddFile(file.Path, &file.Hash, &otherFile.Hash)
				}
				found = true
				break
			}
		}

		if !found {
			diff.AddFile(otherFile.Path, nil, &otherFile.Hash)
		}
	}

	return &diff
}

func (d *Diff) AddFile(path string, reference *uint32, given *uint32) {
	d.Files = append(d.Files, DiffFile{path, reference, given})
}

func (d *Diff) AddName(reference string, given string) {
	if reference != given {
		d.Plugin = &DiffName{&reference, &given}
	}
}

func (d *Diff) AddVersion(reference string, given string) {
	if reference != given {
		d.Version = &DiffName{&reference, &given}
	}
}

func (s *Scan) Scan(path string) {
	if strings.HasSuffix(path, ".php") {
		f, err := os.Open(path)
		if err != nil {
			s.AddErrored(path, err)
			return
		}

		hash, err := GetHash(f)
		if err != nil {
			s.AddErrored(path, err)
			return
		}

		s.AddHashed(path, hash)
	}
}
