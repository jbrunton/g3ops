package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/logrusorgru/aurora"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
)

// WorkflowDefinition - definitoin for a workflow defined by a G3ops template
type WorkflowDefinition struct {
	Name        string
	Source      string
	Destination string
	Content     string
}

func getWorkflowTemplates(fs *afero.Afero, workflowsDir string, context *G3opsContext) []string {
	files := []string{}
	err := fs.Walk(workflowsDir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jsonnet" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func getWorkflowName(workflowsDir string, filename string) string {
	templateDir, templateFileName := filepath.Split(filename)
	if templateFileName == "template.jsonnet" {
		// Check to see if the file is a top level template.
		if filepath.Clean(templateDir) != filepath.Clean(workflowsDir) {
			// If the file is called template.jsonnet and it's in a subdirectory, then rename it to the directory name.
			// E.g. "workflows/my-workflow/template.jsonnet" returns "my-workflow"
			return filepath.Base(templateDir)
		}
	}
	// In all other cases, simply return the name of the file less the extension.
	// E.g. "workflows/my-workflow.jsonnet" returns "my-workflow"
	return strings.TrimSuffix(templateFileName, filepath.Ext(templateFileName))
}

// GetWorkflowDefinitions - get workflow definitions for the given context
func GetWorkflowDefinitions(fs *afero.Afero, context *G3opsContext) []*WorkflowDefinition {
	vm := jsonnet.MakeVM()
	vm.StringOutput = true
	vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)

	workflowsDir := filepath.Join(context.Dir, "/workflows")
	templates := getWorkflowTemplates(fs, workflowsDir, context)
	definitions := []*WorkflowDefinition{}
	for _, templatePath := range templates {
		workflowName := getWorkflowName(workflowsDir, templatePath)
		input, err := fs.ReadFile(templatePath)
		if err != nil {
			panic(err)
		}
		workflow, err := vm.EvaluateSnippet(templatePath, string(input))
		if err != nil {
			fmt.Println(styles.StyleError(fmt.Sprintf("Error processing %s", templatePath)))
			fmt.Println(err)
			os.Exit(1)
		}
		destinationPath := ".github/workflows/" + workflowName + ".yml"
		meta := strings.Join([]string{
			"# File generated by g3ops, do not modify",
			fmt.Sprintf("# Source: %s", templatePath),
		}, "\n")
		definition := &WorkflowDefinition{
			Name:        workflowName,
			Source:      templatePath,
			Destination: destinationPath,
			Content:     meta + "\n" + workflow,
		}
		definitions = append(definitions, definition)
	}

	return definitions
}

// InitWorkflows - copies g3ops workflow sources to context directory
func InitWorkflows(fs *afero.Afero, context *G3opsContext) {
	generator := workflowGenerator{
		name: "g3ops",
		sources: []string{
			"/workflows/common/git.libsonnet",
			"/workflows/g3ops/config.libsonnet",
			"/workflows/g3ops/template.jsonnet",
		},
	}
	applyGenerator(fs, context, generator)
}

// UpdateWorkflows - update workflow files for the given context
func UpdateWorkflows(fs *afero.Afero, context *G3opsContext) {
	definitions := GetWorkflowDefinitions(fs, context)
	for _, definition := range definitions {
		updateFileContent(fs, definition.Destination, definition.Content, fmt.Sprintf("(from %s)", definition.Source))
	}
}

// ValidateWorkflows - returns an error if the workflows are out of date
func ValidateWorkflows(fs *afero.Afero, context *G3opsContext) error {
	WorkflowValidator := NewWorkflowValidator(fs)
	definitions := GetWorkflowDefinitions(fs, context)
	valid := true
	for _, definition := range definitions {
		fmt.Printf("Checking %s ... ", aurora.Bold(definition.Name))

		schemaResult := WorkflowValidator.ValidateSchema(definition)
		if !schemaResult.Valid {
			fmt.Println(styles.StyleError("FAILED"))
			fmt.Println("  Workflow failed schema validation:")
			for _, err := range schemaResult.Errors {
				fmt.Printf("  ► %s\n", err)
			}
			valid = false
			continue
		}

		contentResult := WorkflowValidator.ValidateContent(definition)
		if !contentResult.Valid {
			fmt.Println(styles.StyleError("FAILED"))
			fmt.Println("  " + contentResult.Errors[0])
			fmt.Println("  ► Run \"g3ops workflow generate\" to update")
			valid = false
			continue
		}

		fmt.Println(styles.StyleOK("OK"))
	}
	if !valid {
		return errors.New("workflow validation failed")
	}
	return nil
}
