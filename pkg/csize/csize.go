package csize

import "fmt"

const (
	Byte      = 1
	Kilobyte  = 1024 * Byte
	Megabyte  = 1024 * Kilobyte
	Gigabyte  = 1024 * Megabyte
	Terabyte  = 1024 * Gigabyte
	Zettabyte = 1024 * Terabyte
)

func FormatSize(size int64) string {
	switch {
	case size >= Zettabyte:
		return fmt.Sprintf("%.2f ZB", float64(size)/Zettabyte)
	case size >= Terabyte:
		return fmt.Sprintf("%.2f TB", float64(size)/Terabyte)
	case size >= Gigabyte:
		return fmt.Sprintf("%.2f GB", float64(size)/Gigabyte)
	case size >= Megabyte:
		return fmt.Sprintf("%.2f MB", float64(size)/Megabyte)
	case size >= Kilobyte:
		return fmt.Sprintf("%.2f KB", float64(size)/Kilobyte)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
