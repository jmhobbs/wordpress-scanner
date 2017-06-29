package shared

import (
	"io"

	"github.com/spaolacci/murmur3"
)

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
