package main

import (
	"testing"
)

func TestLuaScript(t *testing.T) {
	s := NewLuaScript("test_script.lua", 5)
	err := s.Start()
	if err != nil {
		t.Error(err)
	}

	// for i := 0; i < 5; i++ {
	s.pool.DChan <- "test"
	// 	s.pool.DChan <- ExecData{
	// 		email: "pinkearwax@email.com",
	// 		pass:  "real!!!123",
	// 	}
	// }

	// for i := 0; i < 5; i++ {
	// 	s.pool.DChan <- ExecData{
	// 		email: "pinkearwax@email.com",
	// 		pass:  "wrongpassword",
	// 	}
	// }

	s.pool.Stop()

	_ = s
}
