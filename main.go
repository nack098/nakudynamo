package main

import (
	"fmt"
	"log"
	"nakuya/nakudynamo/internal"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "nakuya/nakudynamo/internal/launcher"
)

func main() {
	fmt.Println("Starting nakudynamo...")

	dynoEnv, err := internal.PrepareEnvironment()
	if err != nil {
		log.Fatalf("Failed to prepare environment: %v", err)
	}

	process, err := internal.Start(dynoEnv)
	if err != nil {
		log.Fatalf("Failed to start DynamoDB Local: %v", err)
	}
	fmt.Println("DynamoDB Local server started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Println("\n Termination signal received, stopping DynamoDB Local...")
	err = internal.StopDynamoDB(process)
	if err != nil {
		log.Printf("Error stopping DynamoDB Local: %v", err)
	}

	fmt.Println("Done!")
	time.Sleep(time.Second)
}
