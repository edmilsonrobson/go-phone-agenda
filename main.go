package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	portNumber := flag.Int("port", 8000, "a port number where the server will run")
	flag.Parse()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", *portNumber),
	}

	fmt.Printf("Running on http://127.0.0.1:%v", *portNumber)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
