package qgit

import (
	"fmt"
	"os"
)

type OutputWriter interface {
	WriteOutput(key, value string) error
}

type FileOutputWriter struct {
	file string
}

func NewFileOutputWriter(file string) *FileOutputWriter {
	return &FileOutputWriter{file: file}
}

func (f *FileOutputWriter) WriteOutput(key, value string) error {
	outFile, err := os.OpenFile(f.file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := outFile.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		return err
	}
	return nil
}
