package profile

import "testing"

func TestServiceAndFileStoreLifecycle(t *testing.T) {
	service := NewService(NewFileStore(t.TempDir()))

	original := Profile{Name: "test", User: "example", Project: "new-project"}
	if err := service.Create(original); err != nil {
		t.Fatalf("create profile: %v", err)
	}

	got, err := service.Get("test")
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if got != original {
		t.Fatalf("profile mismatch: got %#v, want %#v", got, original)
	}

	profiles, err := service.List()
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if len(profiles) != 1 || profiles[0] != original {
		t.Fatalf("unexpected profiles: %#v", profiles)
	}

	if err := service.Delete("test"); err != nil {
		t.Fatalf("delete profile: %v", err)
	}
	if _, err := service.Get("test"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestProfileNameValidation(t *testing.T) {
	service := NewService(NewFileStore(t.TempDir()))

	err := service.Create(Profile{Name: "../bad", User: "u", Project: "p"})
	if err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

func TestDuplicateProfile(t *testing.T) {
	service := NewService(NewFileStore(t.TempDir()))

	p := Profile{Name: "test", User: "example", Project: "project"}
	if err := service.Create(p); err != nil {
		t.Fatalf("create profile: %v", err)
	}
	if err := service.Create(p); err != ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}
