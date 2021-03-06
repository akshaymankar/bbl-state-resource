package outrunner

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/cloudfoundry/bbl-state-resource/concourse"
)

// Params.Name > Params.Namefile > Params.StateDir/name > random
func Name(sourcesDir string, params concourse.OutParams) (string, error) {
	if params.Name != "" {
		return params.Name, nil
	}

	nameFilePath := params.NameFile
	if nameFilePath == "" && params.StateDir != "" {
		nameFilePath = filepath.Join(params.StateDir, "name")
	}

	if nameFilePath != "" {
		name, err := ioutil.ReadFile(filepath.Join(sourcesDir, nameFilePath))
		if err != nil {
			return "", fmt.Errorf("Failure reading name file: %s", err)
		}
		return strings.TrimSpace(string(name)), nil
	}

	return fmt.Sprintf("%s-%s", randomdata.Adjective(), randomdata.Noun()), nil
}
