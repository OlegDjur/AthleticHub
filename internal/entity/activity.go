package entity

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"time"
)

// Workout представляет структуру данных для одной тренировки
type Workout struct {
	ID           uuid.UUID     `json:"id" db:"id"`                         // уникальный идентификатор тренировки
	UserID       uuid.UUID     `json:"-" db:"user_id"`                     // id пользователя
	Name         string        `json:"name" db:"name"`                     // Имя тренировки
	Description  string        `json:"description" db:"description"`       // Описание тренировки (заметка)
	Date         time.Time     `json:"date" db:"date"`                     // дата проведения тренировки в формате YYYY-MM-DD
	Duration     time.Duration `json:"duration" db:"duration"`             // продолжительность тренировки (без пауз) в формате ММ:СС или ЧЧ:ММ:СС
	Distance     string        `json:"distance" db:"distance"`             // дистанция тренировки в километрах Форматируется с 2 знаками после запятой (например, "5.47")
	AvgPace      int           `json:"avg_pace" db:"avg_pace"`             // средний темп в формате "М:СС мин/км" Рассчитывается на основе дистанции и активного времени
	AvgHeartRate int           `json:"avg_heart_rate" db:"avg_heart_rate"` // средний пульс за тренировку в ударах в минуту
	MaxHeartRate int           `json:"max_heart_rate" db:"avg_heart_rate"` // средний пульс за тренировку в ударах в минуту
	AvgCadence   uint8         `json:"avg_cadence"`                        // Средний каденс
	SportType    string        `json:"sport_type" db:"sport_type"`         // тип спорта в человеко-читаемом формате Примеры: "Бег", "Велосипед", "Плавание", "Ходьба"
	Calories     uint16        `json:"calories" db:"calories"`             // Каллории
	CreatedAt    time.Time     `json:"-" db:"created_at"`                  // время создания записи в системе Автоматически устанавливается при добавлении тренировки
	UpdatedAt    time.Time     `json:"-" db:"updated_at"`                  // время последнего обновления записи
	RecordData   []RecordData  `json:"record_data" db:"record_data"`
}

func convertDistance(distance uint32) string {
	// Конвертируем дистанцию из сантиметров в километры
	// Делим на 100000 (100 см в метре * 1000 метров в км)
	d := float64(distance) / 100000.0

	return fmt.Sprintf("%.2f", d)
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
	// {"1", "2025-06-12", "10.0", "3600", "6.0", 140},
	// {"2", "2025-06-11", "5.0", "1800", "6.0", 135},
}
