package shared

import (
	"bytes"
	"encoding/binary"
)

// [uint16 bytes][plugin][uint8 bytes][version]
// [uint16 bytes][file path][uint32 size]

func (s *Scan) MarshalToBinary() ([]byte, error) {
	// TODO: Check lengths and fail as needed.
	// OR, switch to variable length encoding.
	msg := new(bytes.Buffer)

	// 1. Write the Scan header
	length := len(s.Plugin)
	err := binary.Write(msg, binary.LittleEndian, uint16(length))
	if err != nil {
		return []byte{}, err
	}

	err = binary.Write(msg, binary.LittleEndian, []byte(s.Plugin))
	if err != nil {
		return []byte{}, err
	}

	length = len(s.Version)
	err = msg.WriteByte(uint8(length))
	if err != nil {
		return []byte{}, err
	}
	err = binary.Write(msg, binary.LittleEndian, []byte(s.Version))
	if err != nil {
		return []byte{}, err
	}

	// 2. Iterate over files, writing them out
	for _, f := range s.Files {
		length = len(f.Path)
		err = binary.Write(msg, binary.LittleEndian, uint16(length))
		if err != nil {
			return []byte{}, err
		}

		err = binary.Write(msg, binary.LittleEndian, []byte(f.Path))
		if err != nil {
			return []byte{}, err
		}

		err = binary.Write(msg, binary.LittleEndian, f.Hash)
		if err != nil {
			return []byte{}, err
		}
	}

	return msg.Bytes(), nil
}
