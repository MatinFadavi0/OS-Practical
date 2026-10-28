package main

import (
	"context"
	"os"
	"os/signal"

	"langar-lab/internal/lab"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	os.Exit(lab.New(os.Stdin, os.Stdout, os.Stderr).Run(ctx, os.Args[1:]))
}
