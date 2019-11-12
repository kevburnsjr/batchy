package main

import (
	"io/ioutil"
	"time"
)

type DataWriter interface {
	Write(data []byte) error
}

type dataWriter struct{}

func (r *dataWriter) Write(data []byte) error {
	return ioutil.WriteFile("test-"+time.Now().String(), data, 0644)
}

func NewDataWriter() *dataWriter {
	return &dataWriter{}
}
