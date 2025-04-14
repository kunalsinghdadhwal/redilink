package routes

import (
	"time"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL                string        `json:"url"`
	CustomShort        string        `json:"short"`
	Expiry             time.Duration `json:"expiry"`
	X_Rate_Remaining   int           `json:"rate_limit"`
	X_Rate_Limit_Reset time.Duration `json:"rate_limit_reset"`
}
