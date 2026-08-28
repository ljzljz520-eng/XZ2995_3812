package cli

import (
	"errors"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(line string) (Command, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Command{}, errors.New("empty command")
	}
	return Command{Name: parts[0], Args: append([]string(nil), parts[1:]...)}, nil
}

func (c Command) Valid() bool {
	switch c.Name {
	case "register":
		return len(c.Args) == 2
	case "issue", "show":
		return len(c.Args) == 1
	case "health":
		return len(c.Args) == 0
	default:
		return false
	}
}

func (c Command) Normalized() Command {
	args := make([]string, len(c.Args))
	copy(args, c.Args)
	for index, arg := range args {
		args[index] = strings.TrimSpace(arg)
	}
	return Command{Name: strings.ToLower(strings.TrimSpace(c.Name)), Args: args}
}

func Commands() []string {
	return []string{"register", "issue", "show", "health"}
}

func IsKnownCommand(name string) bool {
	for _, command := range Commands() {
		if command == name {
			return true
		}
	}
	return false
}

func ParseAndValidate(line string) (Command, error) {
	command, err := Parse(line)
	if err != nil {
		return Command{}, err
	}
	command = command.Normalized()
	if !command.Valid() {
		return Command{}, errors.New("invalid command arguments")
	}
	return command, nil
}
