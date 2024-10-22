package main

import (
	"context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := example(ctx); err != nil {
		log.Println(err)
	}
}
