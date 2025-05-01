package utils

import (
    "fmt"
    "os"
    "sync"
)

// FileWriter is a singleton that writes formatted messages to a file
type FileWriter struct {
    file  *os.File
    mutex sync.Mutex
}

var (
    instance *FileWriter
    once     sync.Once
)

// GetFileWriter returns the singleton instance of FileWriter
func GetFileWriter() *FileWriter {
    once.Do(func() {
        instance = &FileWriter{}
        err := instance.open()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
        }
    })
    return instance
}

// open opens the file
func (fw *FileWriter) open() error {
    // Hardcoded filename
    filename := "application.log"
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    fw.file = file
    return nil
}

// Write writes a formatted message to the file
func (fw *FileWriter) Write(format string, args ...interface{}) error {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    // Reopen the file if it was closed or nil
    if fw.file == nil {
        if err := fw.open(); err != nil {
            return err
        }
    }

    message := fmt.Sprintf(format, args...)
    _, err := fmt.Fprintln(fw.file, message)
    return err
}

// Close closes the file
func (fw *FileWriter) Close() error {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    if fw.file != nil {
        err := fw.file.Close()
        fw.file = nil
        return err
    }
    return nil
}





package main

import (
    "fmt"
    "your-module-name/utils"
)

func main() {
    // Get the FileWriter instance
    writer := utils.GetFileWriter()
    
    // Write to file
    err := writer.Write("Application started")
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }
    
    // Write formatted message
    writer.Write("Processing item %d", 42)
    
    // Make sure to close the file writer when done
    defer utils.GetFileWriter().Close()
    
    // Your main program logic here
}
