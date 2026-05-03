package mws

import (
	"bytes"
	"strings"
	"testing"

	"mws-cli/internal/profile"
)

func newTestApp(t *testing.T) (*App, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := NewWithStore(&out, &errOut, profile.NewFileStore(t.TempDir()))
	return app, &out, &errOut
}

func TestProfileCommands(t *testing.T) {
	app, out, errOut := newTestApp(t)

	if code := app.Run([]string{"profile", "create", "--name", "test", "--user", "example", "--project", "new-project"}); code != 0 {
		t.Fatalf("create exit code = %d, stderr = %s", code, errOut.String())
	}

	out.Reset()
	if code := app.Run([]string{"profile", "get", "test"}); code != 0 {
		t.Fatalf("get exit code = %d, stderr = %s", code, errOut.String())
	}
	if got := out.String(); !strings.Contains(got, "+------+---------+-------------+") || !strings.Contains(got, "| test | example | new-project |") {
		t.Fatalf("unexpected get output:\n%s", got)
	}

	out.Reset()
	if code := app.Run([]string{"profile", "list"}); code != 0 {
		t.Fatalf("list exit code = %d, stderr = %s", code, errOut.String())
	}
	if got := out.String(); !strings.Contains(got, "| NAME | USER    | PROJECT     |") || !strings.Contains(got, "| test | example | new-project |") {
		t.Fatalf("unexpected list output:\n%s", got)
	}

	out.Reset()
	if code := app.Run([]string{"profile", "delete", "test"}); code != 0 {
		t.Fatalf("delete exit code = %d, stderr = %s", code, errOut.String())
	}
}

func TestProfileCreateWithPositionalName(t *testing.T) {
	app, _, errOut := newTestApp(t)

	code := app.Run([]string{"profile", "create", "prod", "--user", "u", "--project", "p"})
	if code != exitUsage {
		t.Fatalf("expected usage exit code %d, got %d", exitUsage, code)
	}
	if got := errOut.String(); !strings.Contains(got, "profile create does not accept positional arguments") {
		t.Fatalf("unexpected stderr:\n%s", got)
	}
}

func TestProfileGetJSON(t *testing.T) {
	app, out, errOut := newTestApp(t)

	if code := app.Run([]string{"profile", "create", "--name", "test", "--user", "example", "--project", "project"}); code != 0 {
		t.Fatalf("create exit code = %d, stderr = %s", code, errOut.String())
	}

	out.Reset()
	if code := app.Run([]string{"profile", "get", "test", "--output", "json"}); code != 0 {
		t.Fatalf("get json exit code = %d, stderr = %s", code, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, `"name": "test"`) || !strings.Contains(got, `"project": "project"`) {
		t.Fatalf("unexpected json output:\n%s", got)
	}
}

func TestUnknownCommandDoesNotPrintHelp(t *testing.T) {
	app, out, errOut := newTestApp(t)

	code := app.Run([]string{"profile", "bad"})
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if out.Len() != 0 {
		t.Fatalf("stdout must be empty, got:\n%s", out.String())
	}
	if got := errOut.String(); !strings.Contains(got, `unknown profile command "bad"`) {
		t.Fatalf("unexpected stderr:\n%s", got)
	}
}

func TestUsageErrorExitCode(t *testing.T) {
	app, _, errOut := newTestApp(t)

	code := app.Run([]string{"profile", "create", "--bad", "value"})
	if code != exitUsage {
		t.Fatalf("expected usage exit code %d, got %d", exitUsage, code)
	}
	if got := errOut.String(); !strings.Contains(got, "unknown flag --bad") {
		t.Fatalf("unexpected stderr:\n%s", got)
	}
}
