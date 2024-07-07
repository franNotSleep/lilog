package types

import (
	"fmt"
	"io"
)

type Payload interface {
  fmt.Stringer
  io.ReaderFrom
  io.WriterTo
  Bytes() []byte
}
