package usecase

import (
	"runtime/debug"
	"slices"
)

type Monitoring struct{}

type CheckHealthOutput struct {
	Revision string
}

func (m *Monitoring) CheckHealth() *CheckHealthOutput {
	revision := "xxxxxxx"
	if info, ok := debug.ReadBuildInfo(); ok {
		if i := slices.IndexFunc(info.Settings, func(s debug.BuildSetting) bool {
			return s.Key == "vcs.revision"
		}); i != -1 {
			revision = info.Settings[i].Value[:len(revision)]
		}
	}
	return &CheckHealthOutput{Revision: revision}
}
