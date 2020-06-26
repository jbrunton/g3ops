package lib

import "testing"

func TestValidateVersionInvalidFormat(t *testing.T) {
	err := validateVersion("myservice", "2.3.2.1")
	expectedError := `Invalid version name "2.3.2.1", must be a semantic version`
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected %q to be %q", err.Error(), expectedError)
	}
}
