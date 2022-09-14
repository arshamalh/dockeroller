package tools

import (
	"github.com/arshamalh/dockeroller/models"
)

type Configs struct {
	Telegram *models.TelegramInfo
}

// Search for a dockeroller.yml file that includs needed configurations and load it.
func LoadYamlConfig() (Configs, error) {
	// TODO: Implement
	return Configs{}, nil
}
