package main

import (
	"log"
	"sync"

	"github.com/yuin/gopher-lua"
)

type JobOptions struct {
	scriptPath string
	dataPath   string
	proxyPath  string
}

// sometimes threads needs to be changed at runtime. how to deal with this dynamically
type Job struct {
	data       []string
	maxThreads int
	workerChan chan any

	//compiled lua script
	script *lua.FunctionProto

	opts JobOptions
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
	for i := 1; i < j.maxThreads; i++ {
		go j.Worker()
	}
	j.Batcher()

	return nil
}

// i think spawning the max number of threads first, then simply killing either just stop using them or killing them as it needs to scale down?
func (j *Job) Batcher() {
	batch := []any{}
	for {
		batch = append(batch, "this is data")
		if len(batch) == j.maxThreads {
			//this is a little weird... but i think it works?
			j.wg.Add(j.maxThreads)
			j.ExecuteBatch(batch)
			j.wg.Wait()
			batch = nil
		}
	}
}

func (j *Job) ExecuteBatch(batch []any) {
	//send data to pool
	//needs a wait group
	for _, element := range batch {
		j.workerChan <- element
	}

	//re-adjust size of threads if needed here?
}

func (j *Job) Worker() {
	//instead of making a state every time, just make it once then copy it for each worker (not a pointer. not thread safe§)
	//there is a pattern in the docs for a lua pool, look into that.
	L := lua.NewState()
	defer L.Close()

	//dont think this library supports proxies, so will most likely have to fork it or make one custom
	//L.PreloadModule("http", NewHttpModule(&http.Client{}).Loader)

	//gopher-lua LState not thread safe. initiate it at start then use it across worker life
	for elem := range j.workerChan {
		L.SetGlobal("data", lua.LString(elem.(string))) // deciding type is just string?
		err := DoCompiledFile(L, j.script)
		if err != nil {

		}
		j.wg.Done()
		//run the lua script and pass the data
	}
}
