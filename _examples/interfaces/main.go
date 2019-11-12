package main

import (
	"time"
)

func main() {
	{
		dw := NewDataWriter()
		dw.Write([]byte("asdf"))
	}
	{
		dw := NewDataWriterBatched(10, 10 * time.Millisecond)
		dw.Write([]byte("asdf"))
	}
}
