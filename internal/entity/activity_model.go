package entity

import (
	"time"

	"github.com/muktihari/fit/profile/typedef"
)

// ActivityData содержит извлеченные данные из FIT-файла
type ActivityData struct {
	// Основные метрики
	Timestamp      time.Time              // Время начала активности
	LocalTimestamp time.Time              // Локальное время
	TotalDistance  uint32                 // Общая дистанция в сантиметрах
	TotalTimerTime uint32                 // Активное время в миллисекундах
	AvgHeartRate   uint8                  // Средний пульс
	MaxHeartRate   uint8                  // Максимальный пульс
	Sport          typedef.Sport          // Тип спорта
	TriggerMethod  typedef.SessionTrigger // Способ запуска тренировки

	// Дополнительные метрики
	TotalCalories uint16 // Общие калории
	AvgSpeed      uint16 // Средняя скорость в м/с * 1000
	MaxSpeed      uint16 // Максимальная скорость
	AvgCadence    uint8  // Средний каденс
	MaxCadence    uint8  // Максимальный каденс
	TotalAscent   uint16 // Общий подъем в метрах
	TotalDescent  uint16 // Общий спуск в метрах
	AvgPower      uint16 // Средняя мощность (для велосипеда)
	MaxPower      uint16 // Максимальная мощность

	MaxSpeedFromRecords uint16

	// GPS данные
	StartPositionLat int32 // Начальная широта (в semicircles)
	StartPositionLon int32 // Начальная долгота (в semicircles)
	EndPositionLat   int32 // Конечная широта
	EndPositionLon   int32 // Конечная долгота

	// Записи (Records) для детального анализа
	Records []RecordData

	// Круги (Laps)
	// Laps []LapData
}

// RecordData содержит данные одной записи (обычно каждую секунду)
type RecordData struct {
	Timestamp   time.Time
	PositionLat int32  // Широта в semicircles
	PositionLon int32  // Долгота в semicircles
	Distance    uint32 // Накопленная дистанция в сантиметрах
	Speed       uint16 // Скорость в м/с * 1000
	HeartRate   uint8  // Пульс
	Cadence     uint8  // Каденс
	Power       uint16 // Мощность
	Altitude    uint16 // Высота в метрах * 5 + 500
	Temperature int8   // Температура
}

// LapData содержит данные о круге/сегменте
type LapData struct {
	Timestamp        time.Time
	StartTime        time.Time
	TotalElapsedTime uint32             // Общее время в миллисекундах
	TotalTimerTime   uint32             // Активное время в миллисекундах
	TotalDistance    uint32             // Дистанция круга в сантиметрах
	AvgSpeed         uint16             // Средняя скорость
	MaxSpeed         uint16             // Максимальная скорость
	AvgHeartRate     uint8              // Средний пульс
	MaxHeartRate     uint8              // Максимальный пульс
	AvgCadence       uint8              // Средний каденс
	MaxCadence       uint8              // Максимальный каденс
	TotalCalories    uint16             // Калории за круг
	LapTrigger       typedef.LapTrigger // Способ завершения круга
}
