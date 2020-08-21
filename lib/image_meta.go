package lib

import (
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type imageMeta struct {
	Image string
}

func readImageMeta(fs *afero.Afero, path string) (imageMeta, error) {
	input, err := fs.ReadFile(path)
	if err != nil {
		return imageMeta{}, err
	}
	return parseImageMeta(input)
}

func parseImageMeta(input []byte) (meta imageMeta, err error) {
	err = yaml.Unmarshal(input, &meta)
	return
}
