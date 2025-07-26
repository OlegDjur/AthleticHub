package models

import "time"

// Workout - основная модель тренировки
type Workout struct {
	ID        string `json:"id"`         // Уникальный идентификатор
	Date      string `json:"date"`       // Дата тренировки (YYYY-MM-DD)
	Distance  string `json:"distance"`   // Дистанция в километрах
	Duration  string `json:"duration"`   // Продолжительность (HH:MM:SS)
	Pace      string `json:"pace"`       // Средний темп (X:XX мин/км)
	HeartRate int    `json:"heart_rate"` // Средний пульс
	// GPS       []PacePoint `json:"gps,omitempty"` // График темпа
}

// PacePoint - точка графика темпа
type PacePoint struct {
	DistanceKm float64   `json:"distance"`   // Дистанция в км
	Pace       float64   `json:"pace"`       // Темп в мин/км
	HeartRate  int       `json:"heart_rate"` // Средний пульс за интервал
	Time       time.Time `json:"time"`       // Время точки
}

// In-memory хранилище
var Workouts = []Workout{
	{"1", "2025-06-12", "10.0", "3600", "6.0", 140},
	{"2", "2025-06-11", "5.0", "1800", "6.0", 135},
}
