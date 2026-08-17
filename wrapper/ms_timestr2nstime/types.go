package ms_timestr2nstime

import (
	"github.com/bclswl0827/mseedio-wasm/utils"
	"github.com/bclswl0827/mseedio-wasm/wrapper"
)

type IMPL_ms_timestr2nstime struct {
	wrapper.BaseDependency
	stringHandler utils.StringHandler
}
