package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/yuin/gopher-lua"
)

type JobOptions struct {
	scriptPath string
	dataPath   string
	proxyPath  string
}

type Job struct {
	data       []string
	maxThreads int
	workerChan chan any

	//compiled lua script
	script *lua.FunctionProto

	opts JobOptions // this is caca poo poo
	wg   sync.WaitGroup
}

func NewJob(opts JobOptions) *Job {

	return &Job{
		opts:       opts,
		wg:         sync.WaitGroup{},
		maxThreads: 5,
		workerChan: make(chan any),
	}
}

func (j Job) Start() error {
	log.Println("Starting job")
	//compile script
	var err error
	j.script, err = CompileLuaFromPath(j.opts.scriptPath)
	if err != nil {
		return err
	}

	//load proxies
	//load data

	//start workers
	resultChans := []chan any{}
	for i := 1; i < j.maxThreads; i++ {
		resultChans = append(resultChans, j.worker())
	}

	//MULTIPLEX RESULTCHANS
	j.handleBatchResults(multiplex(resultChans))
	j.batcher()

	return nil
}

// i think spawning the max number of threads first, then simply killing either just stop using them or killing them as it needs to scale down?
func (j *Job) batcher() {
	batch := []any{}
	for {
		batch = append(batch, "this is data")
		if len(batch) == j.maxThreads {
			//this is a little weird... but i think it works?
			j.wg.Add(j.maxThreads)
			j.executeBatch(batch)
			j.wg.Wait()
			batch = nil
		}
	}
}

func (j *Job) executeBatch(batch []any) {
	//send data to pool
	//needs a wait group
	for _, element := range batch {
		j.workerChan <- element
	}

	//re-adjust size of threads if needed here?
}

func (j *Job) worker() chan any {
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
		for elem := range j.workerChan { //rename to event, and send kill events?
			L.SetGlobal("data", lua.LString(elem.(string))) // deciding type is just string? <update make this LValue
			err := DoCompiledFile(L, j.script)
			if err != nil {

			}

			//get returned response from script (status of script eg: rate limit) <DOESNT ACCOUNT FOR IF SCRIPT RETURNS NOTHING
			retValue := L.Get(-1) // Get the top value on the stack
			L.Pop(1)              // Pop the value from the stack
			fmt.Println(retValue)
			j.wg.Done()
			//run the lua script and pass the data
		}
	}()
	return out
}

// multiplex worker channel responses into 1 then handle here
func (j *Job) handleBatchResults(results chan any) {
	//handle in batches to reduce stress on database

	batch := []any{}

	for result := range results {
		//increment counters

		batch = append(batch, result)

		if len(batch) == 5 {
			//do something with data

			//clear batch
		}

	}
	//handle excess data
	if len(batch) >= 1 {
		//do something with data
	}
}
