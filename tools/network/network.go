package network

import (
	"time"
)

// Config contains configuration for the http, rpc etc...
type Config struct {
	HTTPPort         string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
}
