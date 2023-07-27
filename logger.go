package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetLogger : Returns logger with configuration
func GetLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	return log.With().Caller().Logger()
}
