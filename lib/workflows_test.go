package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func newTestContext() *G3opsContext {
// 	return &G3opsContext{
// 		Config: &G3opsConfig{
// 		},
// 	}
// }

const testTemplate = `
std.manifestYamlDoc({
  greeting: "Hello, World!"
})
`

func TestGetWorkflowName(t *testing.T) {
	assert.Equal(t, "my-workflow-1", getWorkflowName("/workflows", "/workflows/my-workflow-1.jsonnet"))
	assert.Equal(t, "my-workflow-2", getWorkflowName("/workflows", "/workflows/my-workflow-2/template.jsonnet"))
	assert.Equal(t, "my-workflow-3", getWorkflowName("/workflows", "/workflows/workflows/my-workflow-3.jsonnet"))
	assert.Equal(t, "template", getWorkflowName("/workflows", "/workflows/template.jsonnet"))
}

func TestGenerateWorkflowDefinitions(t *testing.T) {
	context := &G3opsContext{
		Path: ".g3ops/config.yml", // TODO: context path should be the directory, not the path of the config
	}
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/workflows/test.jsonnet", []byte(testTemplate), 0644)

	definitions := generateWorkflowDefinitions(fs, context)

	assert.Len(t, definitions, 1)
	assert.Equal(t, ".g3ops/workflows/test.jsonnet", definitions[0].source)
	assert.Equal(t, ".github/workflows/test.yml", definitions[0].destination)
	assert.Contains(t, definitions[0].content, `"greeting": "Hello, World!"`)
}
