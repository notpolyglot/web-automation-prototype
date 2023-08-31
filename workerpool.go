package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/yuin/gopher-lua"
)

type WorkerPoolOptions struct {
	scriptPath string
	dataPath   string
	proxyPath  string

	maxThreads int
}

type WorkerPool struct {
	data       []string
	workerChan chan any

	//compiled lua script
	script *lua.FunctionProto

	opts WorkerPoolOptions
	wg   sync.WaitGroup
}

func NewWorkerPool(opts WorkerPoolOptions) *WorkerPool {

	return &WorkerPool{
		opts:       opts,
		wg:         sync.WaitGroup{},
		workerChan: make(chan any),
	}
}

func (p WorkerPool) Start() error {
	log.Println("Starting job")
	//compile script
	var err error
	p.script, err = CompileLuaFromPath(p.opts.scriptPath)
	if err != nil {
		return err
	}

	//load proxies
	//load data

	//start workers
	resultChans := []chan any{}
	for i := 1; i < p.opts.maxThreads; i++ {
		resultChans = append(resultChans, p.worker())
	}

	//MULTIPLEX RESULTCHANS
	go p.handleBatchResults(multiplex(resultChans))
	p.producer()

	return nil
}

// i think spawning the max number of threads first, then simply killing either just stop using them or killing them as it needs to scale down?
func (p *WorkerPool) producer() {
	batch := []any{}
	for {
		batch = append(batch, "this is data")
		if len(batch) == p.opts.maxThreads {
			//this is a little weird... but i think it works?
			p.wg.Add(p.opts.maxThreads)
			p.executeBatch(batch)
			p.wg.Wait()
			batch = nil
		}
	}
}

func (p *WorkerPool) executeBatch(batch []any) {
	//send data to pool
	//needs a wait group
	for _, element := range batch {
		p.workerChan <- element
	}

	//re-adjust size of threads if needed here?
}

func (p *WorkerPool) worker() chan any {
	out := make(chan any)
	go func() {
		//instead of making a state every time, just make it once then copy it for each worker (not a pointer. not thread safe§)
		//there is a pattern in the docs for a lua pool, look into that.
		L := lua.NewState()
		defer L.Close()

		//i think  can use 1 lua state but have to make a seperate library then import helpers to use channels to send the data

		//dont think this library supports proxies, so will most likely have to fork it or make one custom
		//L.PreloadModule("http", NewHttpModule(&http.Client{}).Loader)

		//gopher-lua LState not thread safe. initiate it at start then use it across worker life
		for elem := range p.workerChan { //rename to event, and send kill events?
			L.SetGlobal("data", lua.LString(elem.(string))) // deciding type is just string? <update make this LValue
			err := DoCompiledFile(L, p.script)
			if err != nil {

			}

			//get returned response from script (status of script eg: rate limit) <DOESNT ACCOUNT FOR IF SCRIPT RETURNS NOTHING
			retValue := L.Get(-1) // Get the top value on the stack
			L.Pop(1)              // Pop the value from the stack
			_ = retValue
			out <- "retvalue"
			p.wg.Done()
			//run the lua script and pass the data
		}
	}()
	return out
}

// multiplex worker channel responses into 1 then handle here
func (p *WorkerPool) handleBatchResults(results chan any) {
	//handle in batches to reduce stress on database

	batch := []any{}

	for result := range results {
		//increment counters
		batch = append(batch, result)

		if len(batch) == 5 {
			fmt.Println("handling batch :D")
			batch = nil
		}

	}
	//handle excess data
	if len(batch) >= 1 {
	}
}
