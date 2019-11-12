package main

import (
	"io/ioutil"
)

type DataWriter interface {
	Write(data []byte) error
}

type dataWriter struct{}

func (r *dataWriter) Write(data []byte) error {
	return ioutil.WriteFile("test1", data, 0644)
}

func NewDataWriter() *dataWriter {
	return &dataWriter{}
}
