package main

import (
	"testing"
)

func TestJob(t *testing.T) {
	j := NewJob(JobOptions{
		scriptPath: "script.lua",
	})
	err := j.Start()
	if err != nil {
		t.Error(err)
	}
}
