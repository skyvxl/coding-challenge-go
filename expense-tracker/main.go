package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"expense-tracker/cli"
	"expense-tracker/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if err := runApp(ctx, stop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runApp(ctx context.Context, stop context.CancelFunc) error {
	path, err := storage.InitStorage()
	if path == "" && err != nil {
		return err
	}
	rootCmd := cli.NewRootCommand(path)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		stop()
		return err
	}
	stop()
	return nil
}
