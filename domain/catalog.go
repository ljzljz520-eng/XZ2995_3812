package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ModelProfile struct {
	Name               string
	Family             string
	HardwareRevision   string
	CertificateSubject string
}

func ModelCatalog() map[string]ModelProfile {
	return map[string]ModelProfile{
		"ALPHA-1": {Name: "ALPHA-1", Family: "alpha", HardwareRevision: "1", CertificateSubject: "manufacturing/alpha"},
		"BETA-2":  {Name: "BETA-2", Family: "beta", HardwareRevision: "2", CertificateSubject: "manufacturing/beta"},
		"GAMMA-3": {Name: "GAMMA-3", Family: "gamma", HardwareRevision: "3", CertificateSubject: "manufacturing/gamma"},
	}
}

func FindModelProfile(model string) (ModelProfile, bool) {
	profile, ok := ModelCatalog()[NormalizeModel(model)]
	return profile, ok
}

func ModelNames() []string {
	catalog := ModelCatalog()
	names := make([]string, 0, len(catalog))
	for name := range catalog {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func SubjectForDevice(device Device) (string, error) {
	profile, ok := FindModelProfile(device.Model)
	if !ok {
		return "", ErrInvalidModel
	}
	if err := ValidateSerial(device.Serial); err != nil {
		return "", err
	}
	return profile.CertificateSubject + "/" + strings.ToLower(device.Serial), nil
}

func DescribeModel(model string) string {
	profile, ok := FindModelProfile(model)
	if !ok {
		return "unknown model"
	}
	return fmt.Sprintf("%s family=%s revision=%s", profile.Name, profile.Family, profile.HardwareRevision)
}

func ModelsInFamily(family string) []ModelProfile {
	profiles := make([]ModelProfile, 0)
	for _, profile := range ModelCatalog() {
		if profile.Family == family {
			profiles = append(profiles, profile)
		}
	}
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return profiles
}

func CatalogValid() bool {
	for name, profile := range ModelCatalog() {
		if name != profile.Name || ValidateModel(profile.Name) != nil {
			return false
		}
		if profile.Family == "" || profile.HardwareRevision == "" || profile.CertificateSubject == "" {
			return false
		}
	}
	return true
}

func ModelSupported(model string) bool {
	_, ok := FindModelProfile(model)
	return ok
}
