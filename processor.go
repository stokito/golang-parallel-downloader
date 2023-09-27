package main

import "fmt"

type PrintingProcessor struct {
}

func (p *PrintingProcessor) ProcessChunk(chunk []byte) {
	fmt.Printf("Process: %s\n", chunk)
}
