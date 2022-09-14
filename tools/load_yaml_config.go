package tools

import (
	"github.com/arshamalh/dockeroller/contracts"
)

type Configs struct {
	TelegramInfo *contracts.Config
}

// Search for a dockeroller.yml file that includs needed configurations and load it.
func LoadYamlConfig() (Configs, error) {
	// TODO: Implement
	return Configs{}, nil
}
