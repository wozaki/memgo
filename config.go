package memgo

import (
	"time"
)

const defaultConnectTimeout = 1 * time.Second
const defaultThresholdByte = 1024 * 1024

type Config struct {
	// This is the maximum amount of time a client will wait for a connection to complete.
	// The default is 1 second.
	// You can't use 0. If 0, 1 second is used.
	ConnectTimeout time.Duration
	CompressThresholdByte int
}

func (config *Config) connectTimeout() time.Duration {
	if config.ConnectTimeout == 0 {
		return defaultConnectTimeout
	}
	return config.ConnectTimeout
}

func (config *Config) compressThresholdByte() int {
	if config.CompressThresholdByte == 0 {
		return defaultThresholdByte
	}
	return config.CompressThresholdByte
}
