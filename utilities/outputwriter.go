package utilities

import (
	"fmt"
	"os"
)

// OutputWriter is an interface that defines the method to write output as key-value pairs.
type OutputWriter interface {
	// WriteOutput writes a key-value pair to the output destination.
	//
	// Parameters:
	//  - key: The key to be written.
	//  - value: The value to be associated with the key.
	//
	// Returns:
	//  - error: Returns an error if writing the output fails, or nil if successful.
	WriteOutput(key, value string) error
}

// FileOutputWriter is a concrete implementation of OutputWriter
// that writes key-value pairs to a file.
type FileOutputWriter struct {
	file string // file is the path of the file where output is written.
}

// NewFileOutputWriter creates and returns a new instance of FileOutputWriter.
//
// Parameters:
//   - file: The file path where the output will be written.
//
// Returns:
//   - *FileOutputWriter: A pointer to the new FileOutputWriter instance.
func NewFileOutputWriter(file string) *FileOutputWriter {
	return &FileOutputWriter{file: file}
}

// WriteOutput writes a key-value pair to the specified file in the format "key=value".
// If the file doesn't exist, it will append to it; if it does exist, it will append the new data.
//
// Parameters:
//   - key: The key to be written.
//   - value: The value associated with the key.
//
// Returns:
//   - error: Returns an error if writing to the file fails, or nil if successful.
func (f *FileOutputWriter) WriteOutput(key, value string) error {
	outFile, err := os.OpenFile(f.file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the key-value pair to the file in "key=value" format.
	if _, err := outFile.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		return err
	}
	return nil
}
