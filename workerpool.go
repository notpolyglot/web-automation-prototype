package main

import (
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// import "sync"

type WorkerPool struct {
	pool *ants.PoolWithFunc

	//Passes data element
	ExecFN func(i interface{}) ExecResponse

	//thread safe moment to access counters
	//only called if response is received & current time is past last print + minrefreshrate
	PrintFN func()

	database Database

	Logger Logger

	//pause and stop
	pChan chan struct{}
	sChan chan struct{}

	DChan chan any
	rChan chan ExecResponse

	counters map[string]int

	rpm *RollingRPM

	// Logger ants.Logger
	lastPrint      time.Time
	minRefreshRate time.Duration

	disableProducer bool

	wg sync.WaitGroup
}

type ExecResponse struct {
	datum any

	Capture []any
	//list of Counters to increment by 1
	Counters []string

	Success bool

	//logs error, increments error, then skips rest
	Error error

	//increments retry by 1, readds element to queue, then skips rest
	Retry bool
}

var defaultCounters = map[string]int{
	"success": 0,
	"error":   0,
	"retry":   0,
	"ban":     0,
	"fail":    0,
}

// allocates pool but doesnt start executing
func NewWorkerPool(exec func(i any) ExecResponse, size int, opts ...OptFunc) *WorkerPool {
	p := &WorkerPool{
		ExecFN:         exec,
		rpm:            NewMovingRPM(10),
		pChan:          make(chan struct{}),
		sChan:          make(chan struct{}),
		DChan:          make(chan any),
		rChan:          make(chan ExecResponse),
		counters:       defaultCounters,
		Logger:         DefaultLogger{},
		minRefreshRate: time.Second,
	}

	p.pool, _ = ants.NewPoolWithFunc(size, p.consumer, ants.WithLogger(p.Logger))
	for _, fn := range opts {
		fn(p)
	}
	return p
}

type OptFunc func(*WorkerPool)

func WithPrintFunc(fn func()) OptFunc {
	return func(p *WorkerPool) {
		p.PrintFN = fn
	}
}

func WithSize(s int) OptFunc {
	return func(p *WorkerPool) {
		p.pool.Tune(s)
	}
}

func WithLogger(l Logger) OptFunc {
	return func(p *WorkerPool) {
		p.Logger = l
	}
}

func WithDatabase(db Database) OptFunc {
	return func(p *WorkerPool) {
		p.database = db
	}
}

// minimum rate at which PrintFN is called
func WithMinRefreshRate(d time.Duration) OptFunc {
	return func(p *WorkerPool) {
		p.minRefreshRate = d
	}
}

func WithDisableDefaultProducer(p *WorkerPool) {
	p.disableProducer = true
}

func WithCounters(counters []string) OptFunc {
	return func(p *WorkerPool) {
		for _, counter := range counters {
			p.counters[counter] = 0
		}
	}
}

func (p *WorkerPool) RPM() float64 {
	return p.rpm.Avg()
}

func (p *WorkerPool) Start() {
	p.sampleRPM()
	p.PrintFN()
	p.lastPrint = time.Now()

	go p.handleBatchReslts()

	if !p.disableProducer {
		go p.producer()
	}

}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) Pause() {
	p.pChan <- struct{}{}
}

func (p *WorkerPool) Stop() {
	p.sChan <- struct{}{}
	p.wg.Wait()
	close(p.rChan)
	//this is kindaaaa weird but i dont know
	p.wg.Add(1)
	p.wg.Wait()
}

func (p *WorkerPool) producer() {
	for {
		select {
		case d := <-p.DChan:
			p.wg.Add(1) //i forgot
			err := p.pool.Invoke(d)
			if err != nil {
				//either a problem with the pool in which case all will fail
				//or problem with the data in which case retrying will do nothing

			}
		case <-p.pChan:
			<-p.pChan
		case <-p.sChan:
			break

		}
	}
}

func (p *WorkerPool) consumer(i interface{}) {
	r := p.ExecFN(i)

	//need to pass the elem so it can be exported with captured data
	r.datum = i
	p.rChan <- r
	p.wg.Done()
}

func (p *WorkerPool) handleBatchReslts() {
	batch := [][]any{}

	defer func() {
		//handle excess data
		if p.database != nil {
			if len(batch) >= 1 {
				p.database.InsertBatch(batch)
			}
		}

		//do a quickie refresh charvaa

		p.sampleRPM()
		p.PrintFN()
		p.wg.Done()
	}()

	for r := range p.rChan {
		//only need to even attempt print if something changed
		if time.Since(p.lastPrint) >= p.minRefreshRate {
			p.sampleRPM()
			p.PrintFN()
			p.lastPrint = time.Now()
		}

		if r.Error != nil {
			p.Logger.Error(r.Error)
			p.counters["error"] += 1
			continue
		}
		if r.Retry {
			p.counters["retry"] += 1
			p.DChan <- r.datum
			continue
		}
		if !r.Success {
			p.counters["fail"] += 1
			continue
		}
		p.counters["success"] += 1

		for _, elem := range r.Counters {
			_, isDefaultCounter := defaultCounters[elem]
			if !isDefaultCounter {
				p.counters[elem] += 1
			}
		}

		if p.database != nil {
			batch = append(batch, append([]any{r.datum}, r.Capture...))
			if len(batch) == 5 {
				p.database.InsertBatch(batch)
				batch = nil
			}
		}
	}
}
