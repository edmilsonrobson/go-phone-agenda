package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/edmilsonrobson/go-phone-agenda/internal/logs"
	"github.com/edmilsonrobson/go-phone-agenda/internal/repositories"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	defaultServerPort := os.Getenv("SERVER_PORT")
	if defaultServerPort == "" {
		log.Fatal("No server port in .env file. Missing .env?")
	}

	defaultServerPortInt, err := strconv.Atoi(defaultServerPort)
	if err != nil {
		log.Fatal(err)
	}

	portNumber := flag.Int("port", defaultServerPortInt, "a port number where the server will run")
	readTimeout := flag.Int("readtimeout", 10, "the timeout (in seconds) for reading")
	writeTimeout := flag.Int("writetimeout", 10, "the timeout (in seconds) for writing")
	flag.Parse()

	// Pass it down to repository -> handlers
	redisConn, err := redis.Dial("tcp", os.Getenv("REDIS_ADDRESS"))
	if err != nil {
		logs.ErrorLogger.Printf("Failed to connect with Redis: %v", err)
	}
	defer redisConn.Close()

	//repository := repositories.NewInMemoryContactRepository()
	repository := repositories.NewRedisContactRepository(&redisConn)

	srv := &http.Server{
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
		Addr:         fmt.Sprintf(":%v", *portNumber),
		Handler:      Routes(repository),
	}

	fmt.Printf("Running on 127.0.0.1:%v\n", *portNumber)
	fmt.Printf("Read timeout: %v seconds | Write timeout: %v seconds\n", *readTimeout, *writeTimeout)
	fmt.Printf("Redis hostname: %v\n", os.Getenv("REDIS_ADDRESS"))
	fmt.Printf("Environment: %v\n", os.Getenv("ENV"))

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
