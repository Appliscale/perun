package cflogger

import "fmt"

type Logger struct {
	errors []string
}

func (logger *Logger) LogError(error string) {
	logger.errors = append(logger.errors, error)
}

func (logger *Logger) LogValidationError(elementName string, error string) {
	logger.errors = append(logger.errors, elementName + ": " + error)
}

func (logger *Logger) PrintErrors() {
	for _, err := range logger.errors {
		fmt.Println(err)
	}
}