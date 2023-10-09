package starlark

import (
	// "log"
	// "go.starlark.net/starlark"
	"ob3-prototype"

	"go.starlark.net/starlark"
)

type EmbeddedStarlarkScript struct {
	pool   *main.WorkerPool
	thread *starlark.Thread
}

func (s EmbeddedStarlarkScript) execute(datum any) main.ExecResponse {
	s.pool.Logger.Error()
	return main.ExecResponse{}
}

func NewEmbeddedStarlarkScript() *EmbeddedStarlarkScript {

	s := &EmbeddedStarlarkScript{}

	s.thread = &starlark.Thread{}

	s.pool = main.NewWorkerPool(s.execute, 20)
	return s
}
