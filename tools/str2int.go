package tools

import (
	"strconv"

	"github.com/arshamalh/dockeroller/log"
)

// We are sure that the received index is conversable to int,
// so we simplify the error handling by just logging it (for worst case)
func Str2Int(index string) int {
	int, err := strconv.Atoi(index)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return int
}
