# Golang downloader with parallel processing sample

Golang interview live coding task:
Download from URL in specified chunks and process in parallel.
Time for the task: 50 minutes.

Build:

    go build -o downloader

Execute tests:

    go test ./...

Run:

    ./downloader https://data.iana.org/TLD/tlds-alpha-by-domain.txt 100

This will download the file with chunks by 100 bytes.
All these chunks will be gathered and printed after downloading.

Please note: the webserver should return a Content-Length header for the file

