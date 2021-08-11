package kafka

import (
	"time"
)

const (
	clusterCreateTimeout = 120 * time.Minute
	clusterUpdateTimeout = 120 * time.Minute
	clusterDeleteTimeout = 120 * time.Minute
)
