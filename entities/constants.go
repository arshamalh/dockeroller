package entities

import "time"

const (
	LEN_IMG_TRIM         = 12
	LEN_CONT_TRIM        = 12
	STATES_PULL_INTERVAL = time.Second
	// Sleeping time, not too much, not so little (under 500 millisecond would be annoying)
	LOGS_PULL_INTERVAL     = time.Millisecond * 500
	CLIENT_TIMEOUT_SECONDS = 30
)
