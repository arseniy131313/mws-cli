package mws

import (
	"fmt"
	"strings"
)

type stringFlag struct {
	Name  string
	Value *string
}

func parseStringFlags(args []string, flags ...stringFlag) ([]string, error) {
	positionals := make([]string, 0)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positionals = append(positionals, arg)
			continue
		}
		if !strings.HasPrefix(arg, "--") {
			return nil, newUsageError(fmt.Sprintf("unknown shorthand flag %q", arg))
		}

		nameValue := strings.TrimPrefix(arg, "--")
		name, value, hasValue := strings.Cut(nameValue, "=")
		if name == "" {
			return nil, newUsageError("empty flag name")
		}

		flag, ok := findStringFlag(name, flags)
		if !ok {
			return nil, newUsageError(fmt.Sprintf("unknown flag --%s", name))
		}

		if !hasValue {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return nil, newUsageError(fmt.Sprintf("flag --%s requires a value", name))
			}
			i++
			value = args[i]
		}

		*flag.Value = value
	}

	return positionals, nil
}

func findStringFlag(name string, flags []stringFlag) (stringFlag, bool) {
	for _, flag := range flags {
		if flag.Name == name {
			return flag, true
		}
	}
	return stringFlag{}, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func isHelp(value string) bool {
	return value == "help" || value == "--help" || value == "-h"
}
