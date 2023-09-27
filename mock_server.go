package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

var mockHttpServer *httptest.Server
var testChunk = []byte("123456789\n")
var testChunks3 = bytes.Repeat(testChunk, 3)
var testChunks20 = bytes.Repeat(testChunk, 20)

func mockWebserverStart() {
	fileServer := &http.ServeMux{}
	// the /test.txt file has size of 200 bytes e.g. 20 test chunks
	fileServer.HandleFunc("/test.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testChunks20)
	})
	mockHttpServer = httptest.NewServer(fileServer)
}

func mockWebserverStop() {
	if mockHttpServer == nil {
		return
	}
	mockHttpServer.Close()
	mockHttpServer = nil
}
