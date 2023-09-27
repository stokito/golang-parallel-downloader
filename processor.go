package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Processor interface {
	ProcessChunk(chunkNumber int, chunk []byte)
}

// PrintingProcessor prints to strerr the received chunk.
// Used in tests.
type PrintingProcessor struct {
}

func (p *PrintingProcessor) ProcessChunk(chunkNumber int, chunk []byte) {
	_, _ = fmt.Fprintf(os.Stderr, "%s", chunk)
}

// BufProcessor collects all chunks.
// Then we can extract the whole file by calling Bytes() method.
// All chunks are stored into a map where a key is a chunk number.
type BufProcessor struct {
	chunks sync.Map
}

func (p *BufProcessor) ProcessChunk(chunkNumber int, chunk []byte) {
	p.chunks.Store(chunkNumber, chunk)
}

// Bytes Get all collected bytes to check for the result
// The bytes are collected in a sequenced order
func (p *BufProcessor) Bytes() []byte {
	// collect all keys
	keys := make([]int, 0, 0)
	p.chunks.Range(func(key, value interface{}) bool {
		keyInt := key.(int)
		keys = append(keys, keyInt)
		return true
	})
	// sort keys
	sort.Ints(keys)
	// concatenate all chunks ordered by the sorted keys
	buf := bytes.NewBuffer(nil)
	for _, k := range keys {
		chunk, _ := p.chunks.Load(k)
		chunkBytes := chunk.([]byte)
		buf.Write(chunkBytes)
	}
	return buf.Bytes()
}
