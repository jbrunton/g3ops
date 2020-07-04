package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setupValidator(templateContent string) (*afero.Afero, *WorkflowValidator, *WorkflowDefinition) {
	fs := CreateMemFs()
	WorkflowDefinition := newTestWorkflowDefinition("test", templateContent)
	validator := newWorkflowValidator(fs)
	return fs, validator, WorkflowDefinition
}

func TestValidateContent(t *testing.T) {
	fs, validator, definition := setupValidator(exampleTemplate)

	fs.WriteFile(definition.destination, []byte(exampleTemplate), 0644)
	result := validator.ValidateContent(definition)

	assert.True(t, result.valid)
	assert.Equal(t, []string{}, result.errors)
}

func TestValidateContentMissing(t *testing.T) {
	_, validator, definition := setupValidator(exampleTemplate)

	result := validator.ValidateContent(definition)

	assert.False(t, result.valid)
	assert.Equal(t, []string{"Workflow missing for \"test\" (expected workflow at .github/workflows/test.yml)"}, result.errors)
}

func TestValidateContentOutOfDate(t *testing.T) {
	fs, validator, definition := setupValidator(exampleTemplate)

	fs.WriteFile(definition.destination, []byte("incorrect content"), 0644)
	result := validator.ValidateContent(definition)

	assert.False(t, result.valid)
	assert.Equal(t, []string{"Content is out of date for \"test\" (.github/workflows/test.yml)"}, result.errors)
}

func TestValidateSchema(t *testing.T) {
	_, validator, definition := setupValidator(exampleWorkflow)

	result := validator.ValidateSchema(definition)

	assert.True(t, result.valid)
	assert.Equal(t, []string{}, result.errors)
}

func TestValidateSchemaMissingField(t *testing.T) {
	_, validator, definition := setupValidator(invalidWorkflow)

	result := validator.ValidateSchema(definition)

	assert.False(t, result.valid)
	assert.Equal(t, []string{"(root): jobs is required"}, result.errors)
}
