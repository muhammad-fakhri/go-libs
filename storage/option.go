package storage

import (
	"context"

	"gocloud.dev/blob"
)

// Option initializes configs
type Option func(c *config)

type config struct {
	ctx           context.Context
	storageSecret string
	storageBucket string

	// GCP specific settings
	gcpWriterOptions *blob.WriterOptions
	gcpReaderOptions *blob.ReaderOptions
}

func GCPStorage(ctx context.Context, storageBucket, storageSecret string) Option {
	return func(c *config) {
		c.ctx = ctx
		c.storageSecret = storageSecret
		c.storageBucket = storageBucket
	}
}

func GCPReaderOptions(gcpReaderOptions *blob.ReaderOptions) Option {
	return func(c *config) {
		c.gcpReaderOptions = gcpReaderOptions
	}
}

func GCPWriterOptions(gcpWriterOptions *blob.WriterOptions) Option {
	return func(c *config) {
		c.gcpWriterOptions = gcpWriterOptions
	}
}
