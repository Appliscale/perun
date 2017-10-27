package cflogger

import "fmt"

type Logger struct {
	errors []string
}

func LogError(logger *Logger, error string) {
	logger.errors = append(logger.errors, error)
}

func LogValidationError(logger *Logger, elementName string, error string) {
	logger.errors = append(logger.errors, elementName + ": " + error)
}

func PrintErrors(logger *Logger) {
	for _, err := range logger.errors {
		fmt.Println(err)
	}
}