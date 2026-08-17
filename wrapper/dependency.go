package wrapper

import (
	"context"
	"sync"

	"github.com/tetratelabs/wazero/api"
)

type BaseDependency struct {
	Mutex       *sync.Mutex
	Ctx         context.Context
	FunctionObj api.Function
}

func NewBaseDependency(ctx context.Context, mutex *sync.Mutex, functionObj api.Function) BaseDependency {
	return BaseDependency{
		Mutex:       mutex,
		Ctx:         ctx,
		FunctionObj: functionObj,
	}
}
