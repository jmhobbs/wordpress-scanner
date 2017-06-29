package shared

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

func (s *Scan) AddHashed(path string, hash uint32) {
	s.Files = append(s.Files, File{path, nil, hash})
}

func (s *Scan) AddErrored(path string, err error) {
	s.Files = append(s.Files, File{path, err, 0})
}
