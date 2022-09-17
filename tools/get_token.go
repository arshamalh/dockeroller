package tools

import (
	"fmt"
	"os"

	"github.com/arshamalh/dockeroller/contracts"
)

func GetToken(config *contracts.Config) (token string, err error) {
	if ct := (*config)["TOKEN"]; ct != nil {
		token = ct.(string)
		if token != "" {
			return token, nil
		}
	}
	if envt := os.Getenv("TOKEN"); envt != "" {
		return envt, nil
	}
	return "", fmt.Errorf("there is no token for telegram")
}
