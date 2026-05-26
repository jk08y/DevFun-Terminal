package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jk08y/gsh/internal/config"
	"github.com/jk08y/gsh/internal/shell"
)

const version = "1.0.0"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gsh: config warning: %v\n", err)
		cfg = config.Default()
	}

	sh := shell.New(cfg, version)

	// Graceful shutdown on SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		sh.Cleanup()
		os.Exit(0)
	}()

	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "gsh: %v\n", err)
		os.Exit(1)
	}
}
