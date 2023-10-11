package main

import "flag"

func main() {
	luaPath := flag.String("lua", "", "path of lua script")
	threads := flag.Int("t", 20, "number of threads")

	flag.Parse()

	if luaPath != nil {
		s := NewLuaScript(*luaPath, *threads)
		s.Start()
		s.pool.Wait()
	}

}
