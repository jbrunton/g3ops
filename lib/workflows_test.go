package lib

import (
	"testing"

	_ "github.com/jbrunton/g3ops/statik"
	"github.com/stretchr/testify/assert"
)

const exampleTemplate = `
std.manifestYamlDoc({
  greeting: "Hello, World!"
})
`

const exampleWorkflow = `# File generated by g3ops, do not modify
# Source: .g3ops/workflows/test.jsonnet
"greeting": "Hello, World!"
`

func TestGetWorkflowName(t *testing.T) {
	assert.Equal(t, "my-workflow-1", getWorkflowName("/workflows", "/workflows/my-workflow-1.jsonnet"))
	assert.Equal(t, "my-workflow-2", getWorkflowName("/workflows", "/workflows/my-workflow-2/template.jsonnet"))
	assert.Equal(t, "my-workflow-3", getWorkflowName("/workflows", "/workflows/workflows/my-workflow-3.jsonnet"))
	assert.Equal(t, "template", getWorkflowName("/workflows", "/workflows/template.jsonnet"))
}

func TestGenerateWorkflowDefinitions(t *testing.T) {
	context := &G3opsContext{
		Dir: ".g3ops",
	}
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/workflows/test.jsonnet", []byte(exampleTemplate), 0644)

	definitions := generateWorkflowDefinitions(fs, context)

	assert.Len(t, definitions, 1)
	assert.Equal(t, ".g3ops/workflows/test.jsonnet", definitions[0].source)
	assert.Equal(t, ".github/workflows/test.yml", definitions[0].destination)
	assert.Equal(t, definitions[0].content, exampleWorkflow)
}

func TestValidateWorkflows(t *testing.T) {
	context := &G3opsContext{
		Dir: ".g3ops",
	}
	fs := CreateMemFs()

	fs.WriteFile(".g3ops/workflows/test.jsonnet", []byte(exampleTemplate), 0644)
	err := ValidateWorkflows(fs, context)
	assert.EqualError(t, err, "workflows out of date. Run \"g3ops workflow generate\" to update")

	fs.WriteFile(".github/workflows/test.yml", []byte("incorrect content"), 0644)
	err = ValidateWorkflows(fs, context)
	assert.EqualError(t, err, "workflows out of date. Run \"g3ops workflow generate\" to update")

	fs.WriteFile(".github/workflows/test.yml", []byte(exampleWorkflow), 0644)
	err = ValidateWorkflows(fs, context)
	assert.NoError(t, err)
}

func ExampleValidateWorkflows() {
	context := &G3opsContext{
		Dir: ".g3ops",
	}
	fs := CreateMemFs()

	fs.WriteFile(".g3ops/workflows/test.jsonnet", []byte(exampleTemplate), 0644)
	ValidateWorkflows(fs, context)

	fs.WriteFile(".github/workflows/test.yml", []byte("incorrect content"), 0644)
	ValidateWorkflows(fs, context)

	fs.WriteFile(".github/workflows/test.yml", []byte(exampleWorkflow), 0644)
	ValidateWorkflows(fs, context)

	// Output:
	// Checking test ... [1;31mFAILED[0m
	//   Workflow missing for "test" (expected workflow at .github/workflows/test.yml)
	// Checking test ... [1;31mFAILED[0m
	//   Content is out of date for "test" (at .github/workflows/test.yml)
	// Checking test ... [1;32mOK[0m
}

func ExampleInitWorkflows() {
	context := &G3opsContext{
		Dir: ".g3ops",
	}
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/workflows/common/git.libsonnet", []byte(""), 0644)

	InitWorkflows(fs, context)

	// Output:
	// update .g3ops/workflows/common/git.libsonnet
	// create .g3ops/workflows/g3ops/config.libsonnet
	// create .g3ops/workflows/g3ops/template.jsonnet
}
