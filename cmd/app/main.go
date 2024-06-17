package main

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/server"
	"os"
)

func main() {
	ctx := context.Background()
	srv, err := server.InitializeServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing server: %s\n", err)
		os.Exit(1)
	}
	err = srv.Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %s\n", err)
		os.Exit(1)
	}
}
