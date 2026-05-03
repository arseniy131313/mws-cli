package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	filePerm = 0o600
	dirPerm  = 0o700
)

type FileStore struct {
	dir string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{dir: dir}
}

func DefaultDir() string {
	return "profiles"
}

func (s *FileStore) Create(p Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(s.dir, dirPerm); err != nil {
		return fmt.Errorf("create profile directory: %w", err)
	}

	path := s.path(p.Name)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
	if errors.Is(err, os.ErrExist) {
		return ErrAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create profile file: %w", err)
	}
	defer file.Close()

	data, err := encodeProfile(p)
	if err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return err
	}

	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write profile: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("sync profile: %w", err)
	}

	return nil
}

func (s *FileStore) Get(name string) (Profile, error) {
	if err := ValidateName(name); err != nil {
		return Profile{}, err
	}

	data, err := os.ReadFile(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, ErrNotFound
	}
	if err != nil {
		return Profile{}, fmt.Errorf("read profile: %w", err)
	}

	p, err := decodeProfile(data)
	if err != nil {
		return Profile{}, fmt.Errorf("parse profile %q: %w", name, err)
	}
	p.Name = name

	if strings.TrimSpace(p.User) == "" || strings.TrimSpace(p.Project) == "" {
		return Profile{}, fmt.Errorf("profile %q is corrupted: user and project are required", name)
	}
	return p, nil
}

func (s *FileStore) List() ([]Profile, error) {
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}

	profiles := make([]Profile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		if !profileNameRE.MatchString(name) {
			continue
		}

		p, err := s.Get(name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})
	return profiles, nil
}

func (s *FileStore) Delete(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	if err := os.Remove(s.path(name)); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	return nil
}

func (s *FileStore) path(name string) string {
	return filepath.Join(s.dir, name+".yaml")
}
