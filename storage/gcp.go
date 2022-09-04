package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2/google"

	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
)

const (
	APISCOPE = "https://www.googleapis.com/auth/cloud-platform"
)

type gcpStorage struct {
	ctx    context.Context
	bucket *blob.Bucket
}

func newGCPStorage(cfg *config) (*gcpStorage, error) {
	jsonFile, err := os.Open(cfg.storageSecret)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	creds, err := google.CredentialsFromJSON(cfg.ctx, jsonData, APISCOPE)
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}

	bucket, err := gcsblob.OpenBucket(cfg.ctx, client, cfg.storageBucket, nil)
	if err != nil {
		return nil, err
	}

	return &gcpStorage{
		ctx:    cfg.ctx,
		bucket: bucket,
	}, nil
}

func (m *gcpStorage) Read(filepath string, options ...Option) (io.ReadCloser, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	reader, err := m.bucket.NewReader(m.ctx, filepath, c.gcpReaderOptions)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (m *gcpStorage) Close() error {
	return m.bucket.Close()
}

func (m *gcpStorage) Write(filepath string, options ...Option) (io.WriteCloser, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	writer, err := m.bucket.NewWriter(m.ctx, filepath, c.gcpWriterOptions)

	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (m *gcpStorage) Exists(filepath string) bool {
	exists, _ := m.bucket.Exists(m.ctx, filepath)
	return exists
}
