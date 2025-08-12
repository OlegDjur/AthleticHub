package dto

import (
	"github.com/gofrs/uuid/v5"
	"time"
	"workout/internal/entity"
	"workout/internal/utils"
)

type WorkoutDTO struct {
	ID           uuid.UUID `json:"id" db:"id"`                         // уникальный идентификатор тренировки
	UserID       uuid.UUID `json:"-" db:"user_id"`                     // id пользователя
	Name         string    `json:"name" db:"name"`                     // Имя тренировки
	Description  string    `json:"description" db:"description"`       // Описание тренировки (заметка)
	Date         string    `json:"date" db:"date"`                     // дата проведения тренировки в формате YYYY-MM-DD
	Duration     string    `json:"duration" db:"duration"`             // продолжительность тренировки (без пауз) в формате ММ:СС или ЧЧ:ММ:СС
	Distance     string    `json:"distance" db:"distance"`             // дистанция тренировки в километрах Форматируется с 2 знаками после запятой (например, "5.47")
	AvgPace      string    `json:"avg_pace" db:"avg_pace"`             // средний темп в формате "М:СС мин/км" Рассчитывается на основе дистанции и активного времени
	AvgHeartRate int       `json:"avg_heart_rate" db:"avg_heart_rate"` // средний пульс за тренировку в ударах в минуту
	MaxHeartRate int       `json:"max_heart_rate" db:"avg_heart_rate"` // средний пульс за тренировку в ударах в минуту
	AvgCadence   uint8     `json:"avg_cadence"`                        // Средний каденс
	SportType    string    `json:"sport_type" db:"sport_type"`         // тип спорта в человеко-читаемом формате Примеры: "Бег", "Велосипед", "Плавание", "Ходьба"
	Calories     uint16    `json:"calories" db:"calories"`             // Каллории
	CreatedAt    time.Time `json:"-" db:"created_at"`                  // время создания записи в системе Автоматически устанавливается при добавлении тренировки
	UpdatedAt    time.Time `json:"-" db:"updated_at"`                  // время последнего обновления записи
}

func NewWorkoutDTO(data *entity.ActivityData) *WorkoutDTO {
	distance := convertDistance(data.TotalDistance)
	// Вычисляем темп бега на основе дистанции и времени
	pace := utils.CalculatePace(data.TotalDistance, data.TotalTimerTime)

	date := data.LocalTimestamp.Format("2006-01-02")

	workout := &WorkoutDTO{
		Date:         date,                                                // Используем отформатированную дату
		Distance:     distance,                                            // Форматируем дистанцию с 2 знаками после запятой в строку
		Duration:     utils.SecondsToHMS(int(data.TotalTimerTime) / 1000), // Конвертируем время из миллисекунд в секунды и форматируем как ЧЧ:ММ:СС
		AvgPace:      utils.FormatPace(pace),                              // Форматируем ср темп в читаемый вид (например, "5:30 мин/км")
		AvgHeartRate: int(data.AvgHeartRate),                              // Средний пульс за тренировку, приводим к целому числу
		MaxHeartRate: int(data.MaxHeartRate),                              // Максимальный пульс за тренировку, приводим к целому числу
		AvgCadence:   data.AvgCadence,                                     // Средний каденс
		SportType:    utils.GetSportName(data.Sport),                      // Добавляем информацию о типе спорта и способе запуска
		Calories:     data.TotalCalories,
	}

	return workout
}

func WorkoutMapper(w WorkoutDTO) (*entity.Workout, error) {
	// Конвертируем темп
	pace, err := utils.ParsePaceMMSS(w.AvgPace)
	if err != nil {
		return nil, err
	}
	d, err := time.Parse(time.DateOnly, w.Date)
	if err != nil {
		return nil, err
	}

	t, err := ParseHHMMSS(w.Duration)
	if err != nil {
		return nil, err
	}

	return &entity.Workout{
		UserID:       w.UserID,
		Name:         w.Name,
		Description:  w.Description,
		Date:         d,
		Duration:     t,
		Distance:     w.Distance,
		AvgPace:      pace,
		AvgHeartRate: w.AvgHeartRate,
		MaxHeartRate: w.MaxHeartRate,
		AvgCadence:   w.AvgCadence,
		SportType:    w.SportType,
		Calories:     w.Calories,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func ParseHHMMSS(s string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}
	// t – это дата 0000-01-01 + время, поэтому просто конвертируем в Duration
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second, nil
}

type UpdateWorkout struct {
	ID           uuid.UUID
	Name         *string
	SportType    *string
	Duration     *int32
	Distance     *float32
	AvgPace      *float32
	AvgHeartRate *int16
	MaxHeartRate *int16
	AvgCadence   *int16
	Calories     *int32
	Description  *string
	DeviceName   *string
}
