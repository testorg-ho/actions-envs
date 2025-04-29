package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func periodicPrinter(done <-chan struct{}) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("Hello")
        case <-done:
            fmt.Println("Stopping periodic printer...")
            return
        }
    }
}

func main() {
    // Channel to signal shutdown
    done := make(chan struct{})
    
    // Start the periodic printer
    go periodicPrinter(done)

    // Your main program logic goes here
    fmt.Println("Main program started")
    
    // Wait for interrupt signal to gracefully shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    
    // Signal all goroutines to stop
    fmt.Println("Shutting down...")
    close(done)
    
    // Give goroutines time to clean up
    time.Sleep(1 * time.Second)
    fmt.Println("Program exited")
}
