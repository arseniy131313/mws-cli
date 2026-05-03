package profile

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrRequiredName    = errors.New("profile name is required")
	ErrRequiredUser    = errors.New("profile user is required")
	ErrRequiredProject = errors.New("profile project is required")
	ErrInvalidName     = errors.New("profile name is invalid")
	ErrAlreadyExists   = errors.New("profile already exists")
	ErrNotFound        = errors.New("profile not found")
)

var profileNameRE = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

type Profile struct {
	Name    string `json:"name"`
	User    string `json:"user"`
	Project string `json:"project"`
}

func (p Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrRequiredName
	}
	if !profileNameRE.MatchString(p.Name) {
		return ErrInvalidName
	}
	if strings.TrimSpace(p.User) == "" {
		return ErrRequiredUser
	}
	if strings.TrimSpace(p.Project) == "" {
		return ErrRequiredProject
	}
	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrRequiredName
	}
	if !profileNameRE.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
