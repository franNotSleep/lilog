package options

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/frannotsleep/lilog/types"
)

const (
	MaxLigLogOptionsPayloadSize = 1 << 20 // 1 MB
)

var ErrMaxLigLogOptionsPayloadSize = errors.New("maximum payload size exceeded.")

func (n ServerName) Bytes() []byte  { return []byte(n) }
func (n ServerName) String() string { return string(n) }

func (sn ServerName) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, ServerNameType)

	if err != nil {
		return 0, nil
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(sn)))

	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(sn))
	return int64(o) + n, err
}

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

func Decode(opts *LigLogOptions, r io.Reader) (uint8, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ)

	if err != nil {
		return 0, err
	}

	var payload types.Payload
	switch typ {
	case ServerNameType:
		payload = new(ServerName)
		opts.ServerName = payload.(*ServerName)
	default:
		return typ, errors.New("invalid type.")
	}

	_, err = payload.ReadFrom(r)

	if err != nil {
		return typ, err
	}

	return typ, nil
}
