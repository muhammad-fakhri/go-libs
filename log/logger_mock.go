package log

import logrusTest "github.com/sirupsen/logrus/hooks/test"

/**
 * return a Slog object for logging and a Hook object
 *		to capture the logging's output so that it can
 * 		be tested any further to check whether it has
 * 		the right structure and contents or not.
 */
func NewSLoggerWithTestHook(service string) (SLogger, *logrusTest.Hook) {
	entry, logger := getEntryAndLogger(service)
	return &SLog{entry}, logrusTest.NewLocal(logger)
}
