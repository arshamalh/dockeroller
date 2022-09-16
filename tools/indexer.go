package tools

// Find the index of items in the list with index higher than length or less than length
func Indexer(num, length int) int {
	num %= length
	if num < 0 {
		num += length
	}
	return num
}
