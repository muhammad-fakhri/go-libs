package storage

import (
	"errors"
	"io"
)

//go:generate mockgen -destination mock_storage/mock_storage.go -package mock_storage -source init.go

type Storage interface {
	// Write returns a new writer for a specific file on cloud storage
	Write(filepath string, options ...Option) (io.WriteCloser, error)
	// Read reads a file from online storate
	Read(filepath string, options ...Option) (io.ReadCloser, error)
	// Close closes the storage connection
	Close() error
	// Checks if a file exists on the cloud storage
	Exists(filepath string) bool
}

func NewStorage(impl Implementation, options ...Option) (Storage, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	switch impl {
	case GCP:
		return newGCPStorage(c)
	default:
		return nil, errors.New("implementation not found")
	}
}
