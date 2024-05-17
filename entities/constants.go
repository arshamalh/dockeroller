package entities

import "time"

const (
	LEN_IMG_TRIM         int           = 12
	LEN_CONT_TRIM        int           = 12
	LOGS_QUEUE_LEN       int           = 10
	STATES_PULL_INTERVAL time.Duration = time.Second
	// Sleeping time, not too much, not so little (under 500 millisecond would be annoying)
	LOGS_PULL_INTERVAL time.Duration = 500 * time.Millisecond
	CLIENT_TIMEOUT     time.Duration = 30 * time.Second
)
