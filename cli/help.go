package cli

import (
	"fmt"
	"sort"
	"strings"
)

type HelpEntry struct {
	Command     string
	Arguments   string
	Description string
}

func HelpEntries() []HelpEntry {
	return []HelpEntry{
		{Command: "register", Arguments: "SERIAL MODEL", Description: "validate and persist a device registration"},
		{Command: "issue", Arguments: "SERIAL", Description: "generate ECC request data and issue a certificate"},
		{Command: "show", Arguments: "SERIAL", Description: "display recovered device and certificate state"},
		{Command: "health", Arguments: "", Description: "display persisted entity counts"},
	}
}

func RenderHelp() string {
	entries := HelpEntries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Command < entries[j].Command })
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Command
		if entry.Arguments != "" {
			name += " " + entry.Arguments
		}
		lines = append(lines, fmt.Sprintf("%-28s %s", name, entry.Description))
	}
	return strings.Join(lines, "\n")
}

func HelpFor(command string) (string, bool) {
	for _, entry := range HelpEntries() {
		if entry.Command == command {
			name := entry.Command
			if entry.Arguments != "" {
				name += " " + entry.Arguments
			}
			return name + ": " + entry.Description, true
		}
	}
	return "", false
}

func CommandNames() []string {
	entries := HelpEntries()
	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Command)
	}
	sort.Strings(result)
	return result
}

func UsageLine(command string) string {
	entry, ok := HelpFor(command)
	if !ok {
		return "unknown command"
	}
	return "devicecert " + entry
}
