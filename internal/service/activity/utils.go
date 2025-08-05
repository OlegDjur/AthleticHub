package activity

import (
	"fmt"
)

const (
	FIT_UINT16_INVALID = 0xFFFF // 65535 - для скорости, мощности
)

func isValidFitUint16(value uint16) bool {
	return value != 65535 && value > 0
}

// formatPace форматирует темп в строку вида "5:30 мин/км"
func formatPace(pace float64) string {
	if pace <= 0 {
		return "0:00 мин/км"
	}
	minutes := int(pace)
	seconds := int((pace - float64(minutes)) * 60)
	return fmt.Sprintf("%d:%02d мин/км", minutes, seconds)
}
