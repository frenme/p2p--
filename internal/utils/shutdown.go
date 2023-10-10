package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type GracefulShutdown struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	timeout time.Duration
}

func NewGracefulShutdown(timeout time.Duration) *GracefulShutdown {
	ctx, cancel := context.WithCancel(context.Background())
	return &GracefulShutdown{
		ctx:     ctx,
		cancel:  cancel,
		timeout: timeout,
	}
}

func (gs *GracefulShutdown) Context() context.Context {
	return gs.ctx
}

func (gs *GracefulShutdown) AddTask() {
	gs.wg.Add(1)
}

func (gs *GracefulShutdown) TaskDone() {
	gs.wg.Done()
}

func (gs *GracefulShutdown) WaitForSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigCh
	fmt.Println("\nShutdown signal received...")
	gs.cancel()
	
	done := make(chan struct{})
	go func() {
		gs.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		fmt.Println("Graceful shutdown completed")
	case <-time.After(gs.timeout):
		fmt.Println("Shutdown timeout exceeded, forcing exit")
	}
}