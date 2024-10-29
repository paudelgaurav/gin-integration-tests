package constants

import (
	"time"

	"golang.org/x/time/rate"
)

const (
	AuthRequestPerSecondLimit rate.Limit = 3
	AuthRequestBurstSizeLimit            = 10
	AuthClientCleanupInterval            = 3 * time.Minute
)
