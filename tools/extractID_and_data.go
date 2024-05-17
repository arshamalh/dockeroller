package tools

import "strings"

func ExtractIndexAndID(contextData string) (ID string, index int) {
	rawData := strings.Split(contextData, "|")
	ID = rawData[0]
	index = Str2Int(rawData[1])
	return ID, index
}
