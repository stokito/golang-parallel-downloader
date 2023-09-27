package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
	"time"
)

func TestDownloader_Download(t *testing.T) {
	mockWebserverStart()
	defer mockWebserverStop()
	progressCh := make(chan int, 20)
	d := createTestDownloader(progressCh)
	d.Processor = &PrintingProcessor{}

	err := d.Download()
	if err != nil {
		t.Error(err)
		return
	}

	progress := make([]int, 0, 20)
	for p := range progressCh {
		progress = append(progress, p)
	}
	assert.EqualValues(t, []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}, progress)
	assert.Equal(t, int64(200), d.totalBytes)
}

func TestDownloader_processChunk(t *testing.T) {
	mockWebserverStart()
	defer mockWebserverStop()
	progressCh := make(chan int, 1)
	d := createTestDownloader(progressCh)
	d.totalBytes = 200
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go d.processChunk(1, testChunk, wg)
	p := <-progressCh
	assert.Equal(t, 5, p)
}

type PausedReader struct {
	Reader io.Reader
}

func (r *PausedReader) Read(p []byte) (n int, err error) {
	time.Sleep(4 * time.Second)
	n, err = r.Reader.Read(p)
	return
}

func createTestDownloader(progressCh chan int) *Downloader {
	url := mockHttpServer.URL + "/test.txt"
	chunkSize := 10
	ctx := context.Background()
	d := NewDownloader(ctx, nil, url, progressCh, chunkSize)
	return d
}
