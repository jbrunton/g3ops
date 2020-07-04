package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrunton/g3ops/cmd/styles"
	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	statikFs "github.com/rakyll/statik/fs"
	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
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

type fileSource struct {
	source      string
	destination string
	content     string
}

type workflowGenerator struct {
	name    string
	sources []string
}

func applyFileSource(fs *afero.Afero, content *G3opsContext, source fileSource) {
	var action string
	exists, _ := fs.Exists(source.destination)
	if exists {
		actualContent, _ := fs.ReadFile(source.destination)
		if string(actualContent) == source.content {
			action = "  keep"
		} else {
			action = "update"
		}
	} else {
		action = "create"
	}
	fs.WriteFile(source.destination, []byte(source.content), 0644)
	fmt.Println(action, source.destination)
}

func applyGenerator(fs *afero.Afero, context *G3opsContext, generator workflowGenerator) {
	sourceFs, err := statikFs.New()
	if err != nil {
		panic(err)
	}

	for _, sourcePath := range generator.sources {
		file, err := sourceFs.Open(sourcePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		destinationPath := filepath.Join(context.Dir, sourcePath)
		if err != nil {
			panic(err)
		}
		source := fileSource{
			source:      sourcePath,
			destination: destinationPath,
			content:     string(content),
		}
		applyFileSource(fs, context, source)
	}
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

// func getSchema() *jsonschema.Schema {
// 	schemaLoader := gojsonschema.NewReferenceLoader("https://json.schemastore.org/github-workflow")
// 	// response, err := http.Get("https://json.schemastore.org/github-workflow")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer response.Body.Close()
// 	// data, err := ioutil.ReadAll(response.Body)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// schema := &jsonschema.Schema{}
// 	// if err := json.Unmarshal(data, schema); err != nil {
// 	// 	panic("unmarshal schema: " + err.Error())
// 	// }
// 	// return schema
// }

func convertToStringKeysRecursive(value interface{}, keyPrefix string) (interface{}, error) {
	if mapping, ok := value.(map[interface{}]interface{}); ok {
		dict := make(map[string]interface{})
		for key, entry := range mapping {
			str, ok := key.(string)
			if !ok {
				return nil, formatInvalidKeyError(keyPrefix, key)
			}
			var newKeyPrefix string
			if keyPrefix == "" {
				newKeyPrefix = str
			} else {
				newKeyPrefix = fmt.Sprintf("%s.%s", keyPrefix, str)
			}
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			dict[str] = convertedEntry
		}
		return dict, nil
	}
	if list, ok := value.([]interface{}); ok {
		var convertedList []interface{}
		for index, entry := range list {
			newKeyPrefix := fmt.Sprintf("%s[%d]", keyPrefix, index)
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			convertedList = append(convertedList, convertedEntry)
		}
		return convertedList, nil
	}
	return value, nil
}

func formatInvalidKeyError(keyPrefix string, key interface{}) error {
	var location string
	if keyPrefix == "" {
		location = "at top level"
	} else {
		location = fmt.Sprintf("in %s", keyPrefix)
	}
	return fmt.Errorf("Non-string key %s: %#v", location, key)
}

// ValidateWorkflows - returns an error if the workflows are out of date
func ValidateWorkflows(fs *afero.Afero, context *G3opsContext) error {
	schemaLoader := gojsonschema.NewReferenceLoader("https://json.schemastore.org/github-workflow")
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		panic(err)
	}
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
				//yamlData := map[string]interface{}{"type": "string"}
				var yamlData map[interface{}]interface{}
				err = yaml.Unmarshal(data, &yamlData)
				if err != nil {
					panic(err)
				}
				// var jsonData map[string]interface{}
				// jsonData = make(map[string]interface{})
				// for key, value := range yamlData {
				// 	jsonData[fmt.Sprintf("%v", key)] = value
				// }
				//jsonData := map[string]interface{}{"type": map[string]interface{}{"foo": "bar"}}
				jsonData, err := convertToStringKeysRecursive(yamlData, "")
				if err != nil {
					panic(err)
				}
				//fmt.Printf("%#v\n", jsonData)
				loader := gojsonschema.NewGoLoader(jsonData)
				//result := schema.Validate(context.Context, jsonData)
				result, err := schema.Validate(loader)
				if err != nil {
					panic(err)
				}
				if !result.Valid() {
					fmt.Println(styles.StyleError("FAILED"))
					fmt.Println("  Workflow failed schema validation:")
					for _, err := range result.Errors() {
						// Err implements the ResultError interface
						fmt.Printf("  - %s\n", err)
					}
					valid = false
				} else {
					fmt.Println(styles.StyleCommand("OK"))
				}
			} else {
				valid = false
				fmt.Println(styles.StyleError("FAILED"))
				fmt.Printf("  Content is out of date for %q (at %s)\n", definition.name, definition.destination)
				fmt.Println("  Run \"g3ops workflow generate\" to update")
				//fmt.Printf("Content differs for workflow: %s\nRun \"g3ops workflow generate\" to update", definition.destination)
			}
		} else {
			valid = false
			fmt.Println(styles.StyleError("FAILED"))
			fmt.Printf("  Workflow missing for %q (expected workflow at %s)\n", definition.name, definition.destination)
			fmt.Println("  Run \"g3ops workflow generate\" to update")
		}
	}
	if !valid {
		return errors.New("workflows out of date. Run \"g3ops workflow generate\" to update")
	}
	return nil
}
