package job

import (
    "github.com/yuin/gopher-lua"
	"github.com/cjoudrey/gluahttp"
)

// sometimes threads needs to be changed at runtime. how to deal with this dynamically
type Job struct {
	data string[][]
	maxThreads int
	workerChan chan any

	//compiled lua script
	script *lua.FunctionProto 
}

func NewJob() *Job {

	return &Job{}
}

func (j Job) Start() int {
	go func() {
		go j.Batcher()
	}()
}

//i think spawning the max number of threads first, then simply killing either just stop using them or killing them as it needs to scale down?
func (j Job) Batcher() {
	batch = any[][]
	for {
		batch = append(batch, 0)

		if len(batch) == j.maxThreads {
			j.ExecuteBatch(batch)
		}
	}
}

func (j Job) ExecuteBatch(batch any[][]) {
	//send data to pool
	for _, element := range batch {
		j.workerChan <- element
	}

	//re-adjust size of threads if needed here?
}

func (j Job) Worker() {
	//instead of making a state every time, just make it once then copy it for each worker (not a pointer. not thread safe§)
	//there is a pattern in the docs for a lua pool, look into that.
	L := lua.NewState()
	defer L.Close()

	//dont think this library supports proxies, so will most likely have to fork it or make one custom
	L.PreloadModule("http", NewHttpModule(&http.Client{}).Loader)


	//gopher-lua LState not thread safe. initiate it at start then use it across worker life
	for elem := range j.workerChan {
		//check gopher-lua docs there is a way to store the compiled bytecode instead of doing this.
		if err := L.DoString(`print("hello")`); err != nil {
   			panic(err)
		}
		//run the lua script and pass the data
	}
}
