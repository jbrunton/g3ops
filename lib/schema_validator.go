package lib

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type schemaValidator struct {
	schema *gojsonschema.Schema
}

type validationResult struct {
	valid  bool
	errors []string
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

func convertToStringKeysRecursive(value interface{}, keyPrefix string) (interface{}, error) {
	if mapping, ok := value.(map[interface{}]interface{}); ok {
		return convertToStringDictKeysRecursive(mapping, keyPrefix)
	}
	if list, ok := value.([]interface{}); ok {
		return convertToStringListKeysRecursive(list, keyPrefix)
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

func newSchemaValidator() *schemaValidator {
	schemaLoader := gojsonschema.NewReferenceLoader("https://json.schemastore.org/github-workflow")
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		panic(err)
	}
	return &schemaValidator{
		schema: schema,
	}
}

func (validator *schemaValidator) validate(data []byte) validationResult {
	var yamlData map[interface{}]interface{}
	err := yaml.Unmarshal(data, &yamlData)
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
