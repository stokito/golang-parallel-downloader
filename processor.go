package main

import (
	"bytes"
	"fmt"
	"os"
)

type PrintingProcessor struct {
}

func (p *PrintingProcessor) ProcessChunk(chunk []byte) {
	fmt.Fprintf(os.Stderr, "%s", chunk)
}

type BufProcessor struct {
	buf bytes.Buffer
}

func (p *BufProcessor) ProcessChunk(chunk []byte) {
	_, _ = p.buf.Write(chunk)
}

func (p *BufProcessor) Bytes() []byte {
	return p.buf.Bytes()
}
