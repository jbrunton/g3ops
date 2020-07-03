package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrunton/g3ops/cmd/styles"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
)

type workflowDefinition struct {
	name        string
	source      string
	destination string
	content     string
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

func generateWorkflowDefinitions(fs *afero.Afero, context *G3opsContext) []workflowDefinition {
	vm := jsonnet.MakeVM()
	vm.StringOutput = true
	vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)

	workflowsDir := filepath.Join(context.Dir, "/workflows")
	templates := getWorkflowTemplates(fs, workflowsDir, context)
	definitions := []workflowDefinition{}
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
		definition := workflowDefinition{
			name:        workflowName,
			source:      templatePath,
			destination: destinationPath,
			content:     meta + "\n" + workflow,
		}
		definitions = append(definitions, definition)
	}

	return definitions
}

// GenerateWorkflows - generate workflow files for the given context
func GenerateWorkflows(fs *afero.Afero, context *G3opsContext) {
	definitions := generateWorkflowDefinitions(fs, context)
	for _, definition := range definitions {
		err := fs.WriteFile(definition.destination, []byte(definition.content), 0644)
		fmt.Println("Generated", definition.destination, "from", definition.source)
		if err != nil {
			panic(err)
		}
	}
}

// ValidateWorkflows - returns an error if the workflows are out of date
func ValidateWorkflows(fs *afero.Afero, context *G3opsContext) error {
	definitions := generateWorkflowDefinitions(fs, context)
	valid := true
	for _, definition := range definitions {
		fmt.Printf("Checking %s ... ", definition.name)
		exists, err := fs.Exists(definition.destination)
		if err != nil {
			panic(err)
		}
		if exists {
			data, err := fs.ReadFile(definition.destination)
			if err != nil {
				panic(err)
			}
			actualContent := string(data)
			if actualContent == definition.content {
				fmt.Println(styles.StyleCommand("OK"))
			} else {
				valid = false
				fmt.Println(styles.StyleError("FAILED"))
				fmt.Printf("  Content is out of date for %q (at %s)\n", definition.name, definition.destination)
				//fmt.Printf("Content differs for workflow: %s\nRun \"g3ops workflow generate\" to update", definition.destination)
			}
		} else {
			valid = false
			fmt.Println(styles.StyleError("FAILED"))
			fmt.Printf("  Workflow missing for %q (expected workflow at %s)\n", definition.name, definition.destination)
		}
	}
	if !valid {
		return errors.New("workflows out of date. Run \"g3ops workflow generate\" to update")
	}
	return nil
}