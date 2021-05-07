package main

import (
	"github.com/spaceapegames/lambda-burst/api"
	"log"
	"os"
	"strconv"
)

func main() {
	rateLimit := 1000000 //something so high it'll never get hit
	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		log.Printf("[ERROR] invalid rate limit %s", os.Getenv("RATE_LIMIT"))
	}
	server := api.NewServer(
		os.Getenv("LAMBDA_MODE") == "true",
		8080,
		os.Getenv("LAMBDA_BURST_ADDRESS"),
		rateLimit,
		os.Getenv("DISABLE_BURST") == "true",
	)
	server.Serve()
}
