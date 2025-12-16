package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aprikotdev/speckles/internal/cfg"
	"github.com/aprikotdev/speckles/internal/generator"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	start := time.Now()
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Printf("took %s", time.Since(start))
}

func run(ctx context.Context) error {
	if err := generator.GenerateAll(ctx, "./elements", cfg.Default); err != nil {
		return fmt.Errorf("failed to generate all: %w", err)
	}

	return nil
}
