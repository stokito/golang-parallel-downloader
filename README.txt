# Golang downloader with parallel processing example

Build:

    go build -o downloader

Execute tests:

    go test ./...

Run:

    ./downloader https://www.rfc-editor.org/rfc/rfc1543.txt 100

This will download the file with chunks by 100 bytes.
All these chunks will be gathered and printed after downloading.


