package profile

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type profileFile struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

func encodeProfile(p Profile) ([]byte, error) {
	data, err := yaml.Marshal(profileFile{
		User:    p.User,
		Project: p.Project,
	})
	if err != nil {
		return nil, fmt.Errorf("encode profile yaml: %w", err)
	}

	return data, nil
}

func decodeProfile(data []byte) (Profile, error) {
	var file profileFile

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	if err := decoder.Decode(&file); err != nil {
		return Profile{}, fmt.Errorf("decode profile yaml: %w", err)
	}

	return Profile{
		User:    file.User,
		Project: file.Project,
	}, nil
}
