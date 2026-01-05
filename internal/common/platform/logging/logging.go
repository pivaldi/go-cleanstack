package logging

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func GetLogger() Logger {
	if logger == nil {
		panic("Logging service is not initialized.")
	}

	return logger
}
