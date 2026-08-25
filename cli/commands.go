package cli

import (
	"errors"
	"fmt"
	"strings"

	"devicecert/service"
)

type App struct {
	manager *service.Manager
}

func New(manager *service.Manager) *App { return &App{manager: manager} }

func (a *App) Execute(args []string) (string, error) {
	if a == nil || a.manager == nil {
		return "", errors.New("cli unavailable")
	}
	if len(args) == 0 {
		return Usage(), nil
	}
	switch args[0] {
	case "register":
		return a.register(args[1:])
	case "issue":
		return a.issue(args[1:])
	case "show":
		return a.show(args[1:])
	case "health":
		return a.manager.Health(), nil
	default:
		return "", fmt.Errorf("unknown command %s", args[0])
	}
}

func (a *App) register(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("register requires serial and model")
	}
	device, err := a.manager.RegisterDevice(args[0], args[1])
	if err != nil {
		return "", err
	}
	return service.RenderDevice(device), nil
}

func (a *App) issue(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("issue requires serial")
	}
	result, err := a.manager.IssueCertificate(args[0])
	if err != nil {
		return service.RenderCertificate(result.Certificate), err
	}
	return service.RenderCertificate(result.Certificate), nil
}

func (a *App) show(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("show requires serial")
	}
	recovery, err := a.manager.RecoverDevice(args[0])
	if err != nil {
		return "", err
	}
	return service.RenderRecovery(recovery), nil
}

func Usage() string {
	return strings.Join([]string{"devicecert register DC-12345678 MODEL", "devicecert issue DC-12345678", "devicecert show DC-12345678", "devicecert health"}, "\n")
}

func (a *App) ExecuteLine(line string) (string, error) {
	return a.Execute(strings.Fields(line))
}
