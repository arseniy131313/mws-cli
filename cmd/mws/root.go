package mws

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"mws-cli/internal/profile"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const appName = "mws"

type App struct {
	out     io.Writer
	err     io.Writer
	service *profile.Service
	store   profile.Store
}

func New(out, err io.Writer) *App {
	return &App{out: out, err: err}
}

func NewWithStore(out, err io.Writer, store profile.Store) *App {
	return &App{
		out:     out,
		err:     err,
		store:   store,
		service: profile.NewService(store),
	}
}

func (a *App) Run(args []string) int {
	opts, rest, code := a.parseRoot(args)
	if code != exitOK {
		return code
	}

	if opts.help {
		a.printRootHelp()
		return exitOK
	}
	if opts.version {
		a.printVersion()
		return exitOK
	}

	if a.service == nil {
		dir := profile.DefaultDir()
		a.store = profile.NewFileStore(dir)
		a.service = profile.NewService(a.store)
	}

	if len(rest) == 0 {
		a.printRootHelp()
		return exitOK
	}

	switch rest[0] {
	case "profile":
		return a.runProfile(rest[1:])
	case "version":
		a.printVersion()
		return exitOK
	case "help":
		return a.runHelp(rest[1:])
	default:
		return a.failUsage("unknown command %q\n\nRun '%s --help' to see available commands.", rest[0], appName)
	}
}

type rootOptions struct {
	help    bool
	version bool
}

func (a *App) parseRoot(args []string) (rootOptions, []string, int) {
	var opts rootOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--help" || arg == "-h":
			opts.help = true
			return opts, nil, exitOK
		case arg == "--version" || arg == "-v":
			opts.version = true
			return opts, nil, exitOK
		case strings.HasPrefix(arg, "-"):
			return opts, nil, a.failUsage("unknown global flag %q", arg)
		default:
			return opts, args[i:], exitOK
		}
	}
	return opts, nil, exitOK
}

func (a *App) runHelp(args []string) int {
	if len(args) == 0 {
		a.printRootHelp()
		return exitOK
	}
	if args[0] == "profile" {
		if len(args) == 1 {
			a.printProfileHelp()
			return exitOK
		}
		switch args[1] {
		case "create":
			a.printProfileCreateHelp()
		case "get":
			a.printProfileGetHelp()
		case "list":
			a.printProfileListHelp()
		case "delete":
			a.printProfileDeleteHelp()
		default:
			return a.failUsage("unknown help topic %q", strings.Join(args, " "))
		}
		return exitOK
	}
	return a.failUsage("unknown help topic %q", strings.Join(args, " "))
}

func (a *App) printRootHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`MWS CLI

Usage:
  mws [global flags] <command> [flags]

Available Commands:
  profile     Manage profiles
  version     Print version
  help        Show help

Global Flags:
  -h, --help              Show help
  -v, --version           Print version

Examples:
  mws profile create --name dev --user dev-account --project delivery-dev
  mws profile get dev
  mws profile list
  mws profile delete dev

Use "mws <command> --help" for more information.`)+"\n")
}

func (a *App) printVersion() {
	fmt.Fprintf(a.out, "%s %s\ncommit: %s\nbuilt:  %s\n", appName, version, commit, date)
}

func (a *App) failUsage(format string, args ...any) int {
	fmt.Fprintf(a.err, "Error: "+format+"\n", args...)
	return exitUsage
}

func (a *App) failRuntime(format string, args ...any) int {
	fmt.Fprintf(a.err, "Error: "+format+"\n", args...)
	return exitError
}

func (a *App) profileError(err error) int {
	if err == nil {
		return exitOK
	}

	if errors.Is(err, profile.ErrRequiredName) ||
		errors.Is(err, profile.ErrRequiredUser) ||
		errors.Is(err, profile.ErrRequiredProject) ||
		errors.Is(err, profile.ErrInvalidName) {
		return a.failUsage("%v", err)
	}

	return a.failRuntime("%v", err)
}
