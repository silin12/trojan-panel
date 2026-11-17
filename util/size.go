package util

import "fmt"

// FormatByte 把字节数格式化成人类可读的大小
func FormatByte(b int64) string {
	if b < 0 {
		return "不限"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %c", float64(b)/float64(div), "KMGTPE"[exp])
}
