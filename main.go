package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/go-lua"
	"github.com/alitto/pond"
)

func main() {

	// Create a buffered (non-blocking) pool that can scale up to 100 workers
	// and has a buffer capacity of 1000 tasks
	pool := pond.New(100, 1000)
	start := time.Now()
	// Submit 1000 tasks
	for i := 0; i < 100; i++ {
		pool.Submit(func() {
			l := lua.NewState()
			lua.OpenLibraries(l)
			registerHTTP(l)
			if err := lua.DoFile(l, "script.lua"); err != nil {
				log.Println(err)
			}
		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()

	timeElapsed := time.Since(start)
	fmt.Printf("100 tasks took  %s", timeElapsed)
}

func registerHTTP(l *lua.State) {
	l.Register("sendHTTReq", func(l *lua.State) int {
		arg, ok := l.ToString(1)
		if !ok {

		}

		resp, err := http.Get(arg)
		if err != nil {

		}

		b, err := io.ReadAll(resp.Body)
		// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
		if err != nil {
			// log.Fatalln(err)
		}

		l.PushString(string(b))
		return 1
	})
}
