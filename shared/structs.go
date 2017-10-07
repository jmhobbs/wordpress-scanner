package shared

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
	ReferenceName *string  `json:"reference"`
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
