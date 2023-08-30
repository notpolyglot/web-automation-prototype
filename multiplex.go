package main

import (
	"sync"
)

func multiplex(inputs []<-chan any) chan any {
	var wg sync.WaitGroup
	out := make(chan any)
	for i := range inputs {
		wg.Add(1)
		go func(input <-chan any) {
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
