package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setupValidator() (*afero.Afero, *workflowValidator, *workflowDefinition) {
	fs := CreateMemFs()
	workflowDefinition := newTestWorkflowDefinition("test", exampleTemplate)
	validator := newWorkflowValidator(fs)
	return fs, validator, workflowDefinition
}

func TestValidateValidContent(t *testing.T) {
	fs, validator, definition := setupValidator()

	fs.WriteFile(definition.destination, []byte(exampleTemplate), 0644)
	result := validator.validateContent(definition)

	assert.True(t, result.valid)
	assert.Equal(t, []string{}, result.errors)
}

func TestValidateMissingContent(t *testing.T) {
	_, validator, definition := setupValidator()

	result := validator.validateContent(definition)

	assert.False(t, result.valid)
	assert.Equal(t, []string{"Workflow missing for \"test\" (expected workflow at .github/workflows/test.yml)"}, result.errors)
}

func TestValidateIncorrectContent(t *testing.T) {
	fs, validator, definition := setupValidator()

	fs.WriteFile(definition.destination, []byte("incorrect content"), 0644)
	result := validator.validateContent(definition)

	assert.False(t, result.valid)
	assert.Equal(t, []string{"Content is out of date for \"test\" (.github/workflows/test.yml)"}, result.errors)
}
