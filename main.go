package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	url, chunkSize := parseArgs()
	progressCh := make(chan int, 200)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	processor := &BufProcessor{}
	d := NewDownloader(ctx, processor, url, progressCh, chunkSize)

	go performDownload(d)

	fmt.Printf("Download progress:\n")

	for p := range progressCh {
		fmt.Printf("%d %%\n", p)
	}
	fmt.Printf("%s", processor.Bytes())
}

func parseArgs() (string, int) {
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage downloader <file url> <chunk size>\n"+
			"For example:\n"+
			"downloader https://www.rfc-editor.org/rfc/rfc1543.txt 200\n"+
			"The chunk size may be omitted. Its default is 200 bytes\n")
		os.Exit(1)
	}
	url := os.Args[1]
	chunkSize := 200
	if len(os.Args) >= 3 {
		chunkSizeFromArg, err := strconv.Atoi(os.Args[2])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Invalid chunk size: %s\n", err)
			os.Exit(1)
		}
		chunkSize = chunkSizeFromArg
	}
	return url, chunkSize
}

func performDownload(d *Downloader) {
	err := d.Download()
	if err != nil {
		fmt.Printf("Unable to download from %s: %s", d.Url, err)
		os.Exit(2)
		return
	}
}
