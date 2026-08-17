package mseedio

import (
	"context"
	"sync"

	"github.com/tetratelabs/wazero"
)

func New() *MiniSEED {
	ctx := context.Background()

	return &MiniSEED{
		ctx:     ctx,
		runtime: wazero.NewRuntime(ctx),
		mu:      new(sync.Mutex),
	}
}
