package main

import (
	"fmt"
	"log"
	"os"
)

type Database interface {
	InsertBatch(batch [][]any) error

	Migrate() error
}

// makes use of std lib to output to file
type FlatFileDatabase struct {

	//str fmt eg "value=%s value2=%s intval=%d"
	format string

	logger *log.Logger
}

func NewFlatFileDatabase(path string, format string) *FlatFileDatabase {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//dont think it should be 0
	logger := log.New(file, "", 0)

	// example_fmt := "value1=%s value2=%s"
	return &FlatFileDatabase{
		format: format,
		logger: logger,
	}
}

func (db *FlatFileDatabase) InsertBatch(batch [][]any) error {
	str := ""
	for _, elem := range batch {
		str += fmt.Sprintf(db.format+"\n", elem...)
	}
	db.logger.Print(str)
	return nil
}

func (db FlatFileDatabase) Migrate() error {
	return nil
}

type SQLite3 struct{}

func NewSQLite3() *SQLite3 {

	return &SQLite3{}
}

func (db SQLite3) InsertBatch(batch []any) error {
	return nil
}

func (db SQLite3) Migrate() error {
	return nil
}
