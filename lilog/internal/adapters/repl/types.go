package repl

import (
	"context"

	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Adapter struct {
	api    ports.APIPort
	cancel context.CancelFunc
	ctx    context.Context
}

type OP uint16

const (
	RALL OP = iota + 1
	RONE
	EXIT
)
