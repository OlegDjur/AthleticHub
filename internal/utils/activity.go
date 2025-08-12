package utils

import (
	"fmt"
	"github.com/muktihari/fit/profile/typedef"
	"strconv"
	"strings"
)

// SecondsToHMS конвертирует секунды в формат ЧЧ:ММ:СС
func SecondsToHMS(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// CalculatePace вычисляет темп на основе дистанции и времени
// TotalDistance в сантиметрах, TotalTimerTime в миллисекундах
// Возвращает темп в минутах на километр
func CalculatePace(distance, time uint32) float64 {
	if distance == 0 {
		return 0
	}
	distanceKm := float64(distance) / 100000.0     // см -> км
	timeMinutes := float64(time) / (1000.0 * 60.0) // мс -> минуты

	return timeMinutes / distanceKm
}

// "5:43" -> 343
func ParsePaceMMSS(s string) (int, error) {
	s = strings.TrimSpace(s)
	// Допускаем ведущие/хвостовые пробелы, но никаких единиц измерения.
	if s == "" {
		return 0, fmt.Errorf("pace is empty")
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid pace format, want MM:SS")
	}
	m, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	ss, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	switch {
	case err1 != nil || err2 != nil:
		return 0, fmt.Errorf("invalid numbers in pace")
	case m < 0 || ss < 0 || ss > 59:
		return 0, fmt.Errorf("invalid values in pace")
	}
	return m*60 + ss, nil
}

// FormatPace форматирует темп в строку вида "5:30 мин/км"
func FormatPace(pace float64) string {
	if pace <= 0 {
		return "0:00"
	}
	minutes := int(pace)
	seconds := int((pace - float64(minutes)) * 60)
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// 343 -> "5:43"
func FormatPaceMMSS(sec int) string {
	if sec < 0 {
		sec = 0
	}
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

// GetSportName возвращает человеко-читаемое название вида спорта
func GetSportName(sport typedef.Sport) string {
	switch sport {
	case typedef.SportRunning:
		return "Бег"
	case typedef.SportCycling:
		return "Велосипед"
	case typedef.SportSwimming:
		return "Плавание"
	case typedef.SportWalking:
		return "Ходьба"
	case typedef.SportHiking:
		return "Поход"
	case typedef.SportGeneric:
		return "Общая активность"
	default:
		return fmt.Sprintf("Неизвестный спорт (%d)", sport)
	}
}
