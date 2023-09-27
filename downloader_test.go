package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
	"time"
)

func TestDownloader_Download(t *testing.T) {
	mockWebserverStart()
	defer mockWebserverStop()
	d, progressCh := createTestDownloader()

	err := d.Download()
	if err != nil {
		t.Error(err)
		return
	}

	progress := make([]int, 0, 20)
	for p := range progressCh {
		progress = append(progress, p)
		fmt.Printf("Downloaded: %d %%\n", p)
	}
	assert.EqualValues(t, []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}, progress)
	assert.Equal(t, int64(200), d.totalBytes)
}

func TestDownloader_processChunk(t *testing.T) {
	mockWebserverStart()
	defer mockWebserverStop()
	d, progressCh := createTestDownloader()
	d.totalBytes = 200
	wg := &sync.WaitGroup{}
	wg.Add(1)
	d.processChunk(testChunk, wg)
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

func TestDownloader_processChunk_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url := "http://localhost:8081/test.txt"
	progressCh := make(chan int)
	chunkSize := 10
	d := NewDownloader(ctx, nil, url, progressCh, chunkSize)

	d.totalBytes = 200
	r := &PausedReader{bytes.NewReader(testChunks3)}
	body := io.NopCloser(r)
	go d.processResponse(body)

	p5 := <-progressCh
	assert.Equal(t, 5, p5)
	cancel()
	p10 := <-progressCh
	assert.Equal(t, 0, p10)
	p15 := <-progressCh
	assert.Equal(t, 0, p15)
}

func createTestDownloader() (*Downloader, chan int) {
	url := mockHttpServer.URL + "/test.txt"
	progressCh := make(chan int, 20)
	chunkSize := 10
	ctx := context.Background()
	d := NewDownloader(ctx, nil, url, progressCh, chunkSize)
	return d, progressCh
}
