package main

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/wire"
	"os"
)

func main() {
	ctx := context.Background()
	server, err := wire.InitializeServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing server: %s\n", err)
		os.Exit(1)
	}
	err = server.Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %s\n", err)
		os.Exit(1)
	}
}
