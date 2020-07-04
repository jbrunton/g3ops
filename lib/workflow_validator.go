package lib

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type workflowValidator struct {
	fs     *afero.Afero
	schema *gojsonschema.Schema
}

type validationResult struct {
	valid  bool
	errors []string
}

func convertToStringKeysRecursive(value interface{}, keyPrefix string) (interface{}, error) {
	if mapping, ok := value.(map[interface{}]interface{}); ok {
		return convertToStringDictKeysRecursive(mapping, keyPrefix)
	}
	if list, ok := value.([]interface{}); ok {
		return convertToStringListKeysRecursive(list, keyPrefix)
	}
	return value, nil
}

func convertToStringDictKeysRecursive(mapping map[interface{}]interface{}, keyPrefix string) (interface{}, error) {
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

func convertToStringListKeysRecursive(list []interface{}, keyPrefix string) (interface{}, error) {
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

func formatInvalidKeyError(keyPrefix string, key interface{}) error {
	var location string
	if keyPrefix == "" {
		location = "at top level"
	} else {
		location = fmt.Sprintf("in %s", keyPrefix)
	}
	return fmt.Errorf("Non-string key %s: %#v", location, key)
}

func newWorkflowValidator(fs *afero.Afero) *workflowValidator {
	schemaLoader := gojsonschema.NewReferenceLoader("https://json.schemastore.org/github-workflow")
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		panic(err)
	}
	return &workflowValidator{
		fs:     fs,
		schema: schema,
	}
}

func (validator *workflowValidator) validateSchema(definition *workflowDefinition) validationResult {
	var yamlData map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(definition.content), &yamlData)
	if err != nil {
		panic(err)
	}

	jsonData, err := convertToStringKeysRecursive(yamlData, "")
	if err != nil {
		panic(err)
	}

	loader := gojsonschema.NewGoLoader(jsonData)
	result, err := validator.schema.Validate(loader)
	if err != nil {
		panic(err)
	}

	errors := []string{}
	for _, error := range result.Errors() {
		errors = append(errors, error.String())
	}

	return validationResult{
		valid:  result.Valid(),
		errors: errors,
	}
}

func (validator *workflowValidator) validateContent(definition *workflowDefinition) validationResult {
	exists, err := validator.fs.Exists(definition.destination)
	if err != nil {
		panic(err)
	}

	if !exists {
		reason := fmt.Sprintf("Workflow missing for %s (expected workflow at %s)", aurora.Bold(definition.name), definition.destination)
		return validationResult{
			valid:  false,
			errors: []string{reason},
		}
	}

	data, err := validator.fs.ReadFile(definition.destination)
	if err != nil {
		panic(err)
	}

	actualContent := string(data)
	if actualContent != definition.content {
		reason := fmt.Sprintf("Content is out of date for %s (%s)", aurora.Bold(definition.name), definition.destination)
		return validationResult{
			valid:  false,
			errors: []string{reason},
		}
	}

	return validationResult{
		valid:  true,
		errors: []string{},
	}
}
