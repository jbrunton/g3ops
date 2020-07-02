package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWorkflowName(t *testing.T) {
	assert.Equal(t, "my-workflow-1", getWorkflowName("/workflows", "/workflows/my-workflow-1.jsonnet"))
	assert.Equal(t, "my-workflow-2", getWorkflowName("/workflows", "/workflows/my-workflow-2/template.jsonnet"))
	assert.Equal(t, "my-workflow-3", getWorkflowName("/workflows", "/workflows/workflows/my-workflow-3.jsonnet"))
	assert.Equal(t, "template", getWorkflowName("/workflows", "/workflows/template.jsonnet"))
}
