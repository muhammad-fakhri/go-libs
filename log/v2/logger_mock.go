package log

import (
	"time"

	log "github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
)

/**
 * return a Slog object for logging and a Hook object
 *		to capture the logging's output so that it can
 * 		be tested any further to check whether it has
 * 		the right structure and contents or not.
 */
func NewSLoggerWithTestHook(service string) (SLogger, *logrusTest.Hook) {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	entry := log.NewEntry(logger)
	entry = entry.WithField("service", service)
	return &SLog{entry}, logrusTest.NewLocal(logger)
}
