package tools

import "fmt"

func SizeToHumanReadable(size int64) string {
	if size/1e9 > 0 {
		return fmt.Sprintf("%.2f GigaBytes", float32(size)/1e9)
	} else if size/1e6 > 0 {
		return fmt.Sprintf("%.2f MegaBytes", float32(size)/1e6)
	} else if size/1e3 > 0 {
		return fmt.Sprintf("%.2f KiloBytes", float32(size)/1e3)
	}
	return fmt.Sprintf("%d Bytes", size)
}
