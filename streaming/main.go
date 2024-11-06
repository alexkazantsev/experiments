package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func fileHandler(w http.ResponseWriter, req *http.Request) {
	var (
		err error
		f   *os.File
		b   int64
	)

	if f, err = os.OpenFile("./files/sample.pdf", os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		http.Error(w, "create file error", http.StatusInternalServerError)

		return
	}

	defer f.Close()

	if b, err = io.Copy(f, req.Body); err != nil {
		http.Error(w, "copy stream error", http.StatusInternalServerError)
	}

	_, _ = fmt.Fprintf(w, "done, bytes: %d", b)
}

func main() {
	var router = http.NewServeMux()

	router.HandleFunc("POST /", fileHandler)

	fmt.Println("Starting a new http server")

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("http server err: %+v", err)
	}
}
