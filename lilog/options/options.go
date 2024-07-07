package options

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	MaxLigLogOptionsPayloadSize = 1 << 20 // 1 MB
)

var ErrMaxLigLogOptionsPayloadSize = errors.New("maximum payload size exceeded.")

func (n ServerName) Bytes() []byte  { return []byte(n) }
func (n ServerName) String() string { return string(n) }

func (sn *ServerName) ReadFrom(r io.Reader) (int64, error) {
	var n int64 = 1

	var size uint32
	err := binary.Read(r, binary.BigEndian, &size)

	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxLigLogOptionsPayloadSize {
		return n, ErrMaxLigLogOptionsPayloadSize
	}

	buf := make([]byte, size)
	o, err := r.Read(buf)

	if err != nil {
		return n, err
	}

	*sn = ServerName(buf)
	return n + int64(o), nil
}
