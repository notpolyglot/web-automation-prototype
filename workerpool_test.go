package main

import (
	"testing"
)

func TestWorkerPool(t *testing.T) {
	p := NewWorkerPool(WorkerPoolOptions{
		scriptPath: "script.lua",
		maxThreads: 5,
	})
	err := p.Start()
	if err != nil {
		t.Error(err)
	}
}
