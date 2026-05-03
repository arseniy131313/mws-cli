package mws

import (
	"fmt"
	"strings"

	"mws-cli/internal/output"
	"mws-cli/internal/profile"
)

const (
	formatTable = "table"
	formatJSON  = "json"
)

func (a *App) runProfile(args []string) int {
	if len(args) == 0 || isHelp(args[0]) {
		a.printProfileHelp()
		return exitOK
	}

	switch args[0] {
	case "create":
		return a.profileCreate(args[1:])
	case "get":
		return a.profileGet(args[1:])
	case "list":
		return a.profileList(args[1:])
	case "delete":
		return a.profileDelete(args[1:])
	default:
		return a.failUsage("unknown profile command %q\n\nRun '%s profile --help' to see available profile commands.", args[0], appName)
	}
}

func (a *App) profileCreate(args []string) int {
	if hasHelp(args) {
		a.printProfileCreateHelp()
		return exitOK
	}

	var name, user, project string
	positionals, err := parseStringFlags(args,
		stringFlag{Name: "name", Value: &name},
		stringFlag{Name: "user", Value: &user},
		stringFlag{Name: "project", Value: &project},
	)
	if err != nil {
		return a.commandParseError(err, "mws profile create --help")
	}
	if len(positionals) > 0 {
		return a.failUsage("profile create does not accept positional arguments; use --name")
	}

	p := profile.Profile{Name: name, User: user, Project: project}
	if err := a.service.Create(p); err != nil {
		return a.profileError(err)
	}

	fmt.Fprintf(a.out, "Profile %q created.\n", p.Name)
	return exitOK
}

func (a *App) profileGet(args []string) int {
	if hasHelp(args) {
		a.printProfileGetHelp()
		return exitOK
	}

	var name, format string
	positionals, err := parseStringFlags(args,
		stringFlag{Name: "name", Value: &name},
		stringFlag{Name: "output", Value: &format},
	)
	if err != nil {
		return a.commandParseError(err, "mws profile get --help")
	}
	if len(positionals) > 1 {
		return a.failUsage("profile get accepts exactly one profile name")
	}
	if len(positionals) == 1 {
		name = firstNonEmpty(name, positionals[0])
	}
	format = normalizeOutputFormat(format)
	if !validOutputFormat(format) {
		return a.failUsage("unsupported output format %q; allowed values: table, json", format)
	}

	p, err := a.service.Get(name)
	if err != nil {
		return a.profileError(err)
	}

	return a.renderProfiles([]profile.Profile{p}, format, false)
}

func (a *App) profileList(args []string) int {
	if hasHelp(args) {
		a.printProfileListHelp()
		return exitOK
	}

	var format string
	positionals, err := parseStringFlags(args,
		stringFlag{Name: "output", Value: &format},
	)
	if err != nil {
		return a.commandParseError(err, "mws profile list --help")
	}
	if len(positionals) > 0 {
		return a.failUsage("profile list does not accept positional arguments")
	}
	format = normalizeOutputFormat(format)
	if !validOutputFormat(format) {
		return a.failUsage("unsupported output format %q; allowed values: table, json", format)
	}

	profiles, err := a.service.List()
	if err != nil {
		return a.profileError(err)
	}
	if len(profiles) == 0 && format == formatTable {
		fmt.Fprintln(a.out, "No profiles found.")
		return exitOK
	}

	return a.renderProfiles(profiles, format, true)
}

func (a *App) profileDelete(args []string) int {
	if hasHelp(args) {
		a.printProfileDeleteHelp()
		return exitOK
	}

	var name string
	positionals, err := parseStringFlags(args,
		stringFlag{Name: "name", Value: &name},
	)
	if err != nil {
		return a.commandParseError(err, "mws profile delete --help")
	}
	if len(positionals) > 1 {
		return a.failUsage("profile delete accepts exactly one profile name")
	}
	if len(positionals) == 1 {
		name = firstNonEmpty(name, positionals[0])
	}

	if err := a.service.Delete(name); err != nil {
		return a.profileError(err)
	}

	fmt.Fprintf(a.out, "Profile %q deleted.\n", name)
	return exitOK
}

func (a *App) renderProfiles(profiles []profile.Profile, format string, asList bool) int {
	switch format {
	case formatJSON:
		var err error
		if asList {
			err = output.JSON(a.out, profiles)
		} else if len(profiles) == 1 {
			err = output.JSON(a.out, profiles[0])
		} else {
			err = output.JSON(a.out, profiles)
		}
		if err != nil {
			return a.failRuntime("render json: %v", err)
		}
		return exitOK
	default:
		t := output.NewTable(a.out, "NAME", "USER", "PROJECT")
		for _, p := range profiles {
			t.Row(p.Name, p.User, p.Project)
		}
		if err := t.Render(); err != nil {
			return a.failRuntime("render table: %v", err)
		}
		return exitOK
	}
}

func normalizeOutputFormat(value string) string {
	if strings.TrimSpace(value) == "" {
		return formatTable
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func validOutputFormat(value string) bool {
	return value == formatTable || value == formatJSON
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if isHelp(arg) {
			return true
		}
	}
	return false
}

func (a *App) commandParseError(err error, helpCommand string) int {
	if _, ok := err.(usageError); ok {
		return a.failUsage("%v\n\nRun '%s' for usage.", err, helpCommand)
	}
	return a.failRuntime("%v", err)
}

func (a *App) printProfileHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`Manage profiles.

Usage:
  mws profile <command> [flags]

Available Commands:
  create      Create a profile
  get         Show profile details
  list        List profiles
  delete      Delete a profile

Examples:
  mws profile create --name dev --user dev-account --project delivery-dev
  mws profile get dev
  mws profile list
  mws profile delete dev

Use "mws profile <command> --help" for command-specific help.`)+"\n")
}

func (a *App) printProfileCreateHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`Create a profile.

Usage:
  mws profile create --name <name> --user <user> --project <project>

Flags:
  --name string       Profile name
  --user string       User name
  --project string    Project name

Examples:
  mws profile create --name dev --user dev-account --project delivery-dev`)+"\n")
}

func (a *App) printProfileGetHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`Show profile details.

Usage:
  mws profile get <name> [flags]
  mws profile get --name <name> [flags]

Flags:
  --name string      Profile name; kept for backward compatibility
  --output string    Output format: table, json (default: table)

Examples:
  mws profile get dev
  mws profile get dev --output json`)+"\n")
}

func (a *App) printProfileListHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`List profiles.

Usage:
  mws profile list [flags]

Flags:
  --output string    Output format: table, json (default: table)

Examples:
  mws profile list
  mws profile list --output json`)+"\n")
}

func (a *App) printProfileDeleteHelp() {
	fmt.Fprint(a.out, strings.TrimSpace(`Delete a profile.

Usage:
  mws profile delete <name>
  mws profile delete --name <name>

Flags:
  --name string    Profile name; kept for backward compatibility

Examples:
  mws profile delete dev
  mws profile delete --name dev`)+"\n")
}
