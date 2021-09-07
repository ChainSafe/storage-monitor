package main

import (
	"context"
	"log"

	"github.com/chainsafe/storage-monitor/monitoring"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitor, err := monitoring.New()
	if err != nil {
		log.Fatalln(err)
	}

	err = monitor.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
