package mws

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return e.message
}

func newUsageError(message string) error {
	return usageError{message: message}
}
