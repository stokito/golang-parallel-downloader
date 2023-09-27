package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

type Processor interface {
	ProcessChunk(chunk []byte)
}

type Downloader struct {
	Ctx       context.Context
	Processor Processor
	Url       string
	// ProgressCh receives progress of download of the file in percents
	ProgressCh chan int
	// ChunkSize in bytes
	ChunkSize      int
	totalProcessed atomic.Int64
	totalBytes     int64
}

func NewDownloader(ctx context.Context, processor Processor, url string, progressCh chan int, chunkSize int) *Downloader {
	return &Downloader{
		Ctx:        ctx,
		Processor:  processor,
		Url:        url,
		ProgressCh: progressCh,
		ChunkSize:  chunkSize,
	}
}

func (d *Downloader) Download() error {
	req, err := http.NewRequestWithContext(d.Ctx, "GET", d.Url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Incorrect status code: " + strconv.Itoa(resp.StatusCode))
	}
	d.totalBytes = resp.ContentLength
	err = d.processResponse(resp.Body)
	return err
}

func (d *Downloader) processResponse(repBody io.ReadCloser) error {
	reader := bufio.NewReader(repBody)
	wg := &sync.WaitGroup{}
	err := d.readAndProcess(reader, wg)
	wg.Wait()
	close(d.ProgressCh)
	// Explicitly close body
	_ = repBody.Close()
	return err
}

func (d *Downloader) readAndProcess(reader *bufio.Reader, wg *sync.WaitGroup) error {
	for {
		select {
		case <-d.Ctx.Done():
			return d.Ctx.Err()
		default:
			chunk, err := d.readChunk(reader)
			if err != nil {
				return err
			}
			if chunk == nil {
				return nil
			}
			wg.Add(1)
			go d.processChunk(chunk, wg)
		}
	}
}

func (d *Downloader) readChunk(reader *bufio.Reader) ([]byte, error) {
	chunk := make([]byte, d.ChunkSize)
	readCount, err := reader.Read(chunk)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if readCount == 0 {
		return nil, nil
	}
	chunk = chunk[:readCount]
	return chunk, nil
}

func (d *Downloader) processChunk(chunk []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	if d.Processor != nil {
		d.Processor.ProcessChunk(chunk)
	}
	totalProcessed := d.totalProcessed.Add(int64(len(chunk)))
	progressPercent := totalProcessed * 100 / d.totalBytes
	d.ProgressCh <- int(progressPercent)
}
