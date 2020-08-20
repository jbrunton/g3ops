package lib

import (
	"reflect"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
)

const validConfig = `
name: sandbox

environments:
  production:
    manifest: ./manifests/production.yml

services:
  ping:
    manifest: ./services/ping/manifest.yml

ci:
  defaults:
    build:
      env:
        TAG: $BUILD_VERSION-$BUILD_ID
      command: docker build $BUILD_SERVICE
`

func assertDeepEqual(message string, expected interface{}, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		expectedString, err := yaml.Marshal(&expected)
		if err != nil {
			panic(err)
		}
		actualString, err := yaml.Marshal(&actual)
		if err != nil {
			panic(err)
		}
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expectedString), string(actualString), false)
		t.Errorf("%s\nDiff:\n%s", message, dmp.DiffPrettyText(diffs))
	}
}

// func TestParseConfig(t *testing.T) {
// 	cwd, err := os.Getwd()
// 	config, err := parseConfig([]byte(validConfig))
// 	if err != nil {
// 		t.Errorf("Unexpected error: %q", err)
// 	}

// 	expectedConfig := G3opsConfig{
// 		Name: "sandbox",
// 		Environments: map[string]g3opsEnvironmentConfig{
// 			"production": g3opsEnvironmentConfig{
// 				Manifest: filepath.Join(cwd, "./manifests/production.yml"),
// 			},
// 		},
// 		Services: map[string]g3opsServiceConfig{
// 			"ping": g3opsServiceConfig{
// 				Manifest: filepath.Join(cwd, "./services/ping/manifest.yml"),
// 			},
// 		},
// 		Ci: g3opsCiConfig{
// 			Defaults: g3opsCiDefaultsConfig{
// 				Build: g3opsBuildConfig{
// 					Env: map[string]string{
// 						"TAG": "$BUILD_VERSION-$BUILD_ID",
// 					},
// 					Command: "docker build $BUILD_SERVICE",
// 				},
// 			},
// 		},
// 	}

// 	assertDeepEqual("Mismatch in config", *config, expectedConfig, t)
// }
