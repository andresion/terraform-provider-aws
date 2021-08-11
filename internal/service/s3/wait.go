package s3

import (
	"time"
)

const (
	// Maximum amount of time to wait for S3 changes to propagate
	propagationTimeout = 1 * time.Minute
)
