package storage

type Implementation int

const (
	GCP = Implementation(iota) // Google Cloud Platform Storage
	// can add more implementation here, ex: AWS, Aliyun, ...
)
