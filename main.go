package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	portNumber := flag.Int("port", 8000, "a port number where the server will run")
	readTimeout := flag.Int("readtimeout", 10, "the timeout (in seconds) for reading")
	writeTimeout := flag.Int("writetimeout", 10, "the timeout (in seconds) for writing")

	flag.Parse()

	srv := &http.Server{
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
		Addr:         fmt.Sprintf(":%v", *portNumber),
	}

	fmt.Printf("Running on http://127.0.0.1:%v\n", *portNumber)
	fmt.Printf("Read timeout: %v seconds | Write timeout: %v seconds\n", *readTimeout, *writeTimeout)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
