package main

import (
	"sync"
)

func multiplex[t any](inputs []chan t) chan t {
	var wg sync.WaitGroup
	out := make(chan t)
	for i := range inputs {
		wg.Add(1)
		go func(input <-chan t) {
			for val := range input {
				out <- val
			}
			wg.Done()
		}(inputs[i])
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
