package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceValidation_AddValidationError(t *testing.T) {
	resourceValidation := ResourceValidation{ResourceName: "Name", Errors: []string{}}
	resourceValidation.AddValidationError("Error")
	assert.NotEmpty(t, resourceValidation.Errors)
	assert.Equal(t, "Error", resourceValidation.Errors[0])
}

func TestLogger_HasValidationErrors(t *testing.T) {
	logger := CreateQuietLogger()
	assert.False(t, logger.HasValidationErrors())

	resourceValidation := logger.AddResourceForValidation("Name")
	assert.False(t, logger.HasValidationErrors())

	resourceValidation.AddValidationError("Error")
	assert.True(t, logger.HasValidationErrors())
}

func TestLogger_AddResourceForValidation(t *testing.T) {
	logger := CreateQuietLogger()
	assert.Empty(t, logger.resourceValidation)
	logger.AddResourceForValidation("Name")
	assert.NotEmpty(t, logger.resourceValidation)
}

func TestLogger_SetVerbosity(t *testing.T) {
	logger := CreateQuietLogger()
	logger.SetVerbosity("error")
	assert.Equal(t, ERROR, logger.Verbosity)

	logger.SetVerbosity("Trolololo")
	assert.Equal(t, ERROR, logger.Verbosity)

	logger.SetVerbosity("Warning")
	assert.Equal(t, WARNING, logger.Verbosity)
}

func TestIsVerbosityValid(t *testing.T) {
	assert.True(t, IsVerbosityValid("INFO"))
	assert.True(t, IsVerbosityValid("ERROR"))
	assert.True(t, IsVerbosityValid("WARNING"))
	assert.True(t, IsVerbosityValid("TRACE"))
	assert.True(t, IsVerbosityValid("DEBUG"))

	assert.False(t, IsVerbosityValid("error"))
	assert.False(t, IsVerbosityValid("debug"))
	assert.False(t, IsVerbosityValid("VERBOSE"))

}
