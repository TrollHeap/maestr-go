package logic

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%d min %d s", minutes, seconds)
	}
	return fmt.Sprintf("%d s", seconds)
}
