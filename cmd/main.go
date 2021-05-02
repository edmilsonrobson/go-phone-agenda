package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/edmilsonrobson/go-phone-agenda/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	defaultPortNumber, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	portNumber := flag.Int("port", defaultPortNumber, "a port number where the server will run")
	readTimeout := flag.Int("readtimeout", 10, "the timeout (in seconds) for reading")
	writeTimeout := flag.Int("writetimeout", 10, "the timeout (in seconds) for writing")
	flag.Parse()

	srv := &http.Server{
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
		Addr:         fmt.Sprintf("127.0.0.1:%v", *portNumber),
		Handler:      handlers.Routes(),
	}

	fmt.Printf("Running on http://127.0.0.1:%v\n", *portNumber)
	fmt.Printf("Read timeout: %v seconds | Write timeout: %v seconds\n", *readTimeout, *writeTimeout)
	fmt.Printf("Redis on port %v\n", os.Getenv("REDIS_PORT"))
	fmt.Printf("Environment: %v\n", os.Getenv("ENV"))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
