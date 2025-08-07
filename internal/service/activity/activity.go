package activity

import (
	"bytes"
	"context"
	"fmt"
	"workout/internal/dto"
	"workout/internal/entity"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

type WorkoutService struct {
	// Здесь будет репозиторий для работы с базой данных
	Activity Activity
}

func NewWorkoutService(act Activity) *WorkoutService {
	return &WorkoutService{Activity: act}
}

func (s *WorkoutService) UploadFile(ctx context.Context, data []byte) (*dto.UploadFile, error) {
	// Создаем декодер для библиотеки muktihari/fit
	// Декодер позволяет парсить FIT-файлы с различными опциями
	dec := decoder.New(bytes.NewReader(data))

	// Декодируем FIT-файл и получаем набор протокольных сообщений
	// В muktihari/fit декодирование возвращает массив сообщений proto.Message
	// а не готовую структуру активности как в старой библиотеке
	fit, err := dec.Decode()
	if err != nil {
		// Если файл не является валидным FIT-файлом или произошла ошибка парсинга
		// Возвращаем ошибку 400 "Bad Request" через Echo
		return nil, fmt.Errorf("Ошибка при парсинге FIT-файла: %v", err)
	}

	// Анализируем декодированные сообщения и извлекаем данные активности
	// В новой библиотеке нужно самостоятельно обрабатывать массив сообщений
	activityData, err := extractActivityData(fit.Messages)
	if err != nil {
		return nil, err //c.String(http.StatusInternalServerError, fmt.Sprintf("Ошибка при извлечении данных активности: %v", err))
	}

	// Формируем объект тренировки для сохранения
	newWorkout := dto.NewUploadFile(activityData)

	return newWorkout, nil
}

// extractActivityData извлекает данные активности из массива протокольных сообщений
// В muktihari/fit нужно вручную обрабатывать сообщения для получения данных
func extractActivityData(messages []proto.Message) (*entity.ActivityData, error) {
	// Инициализируем структуру для хранения данных
	data := &entity.ActivityData{}

	// Флаги для отслеживания найденных сообщений
	foundSession := false
	foundActivity := false

	// Перебираем все декодированные сообщения
	for _, message := range messages {
		// Проверяем тип сообщения по глобальному номеру
		switch message.Num {
		case mesgnum.Activity:
			// Высокоуровневый API: автоматическое преобразование в типизированную структуру
			activityMsg := mesgdef.NewActivity(&message)

			if !activityMsg.Timestamp.IsZero() {
				data.Timestamp = activityMsg.Timestamp
			}

			if !activityMsg.LocalTimestamp.IsZero() {
				data.LocalTimestamp = activityMsg.LocalTimestamp
			}

			foundActivity = true
		case mesgnum.Session:
			// Извлекаем основные метрики из первой найденной сессии
			// Высокоуровневый API: готовая структура с типизированными полями
			sessionMsg := mesgdef.NewSession(&message)

			// Типизированные поля с автоматической проверкой на nil
			data.TotalDistance = sessionMsg.TotalDistance
			data.TotalTimerTime = sessionMsg.TotalTimerTime
			data.AvgHeartRate = sessionMsg.AvgHeartRate
			data.MaxHeartRate = sessionMsg.MaxHeartRate
			data.Sport = sessionMsg.Sport
			data.TotalCalories = sessionMsg.TotalCalories
			data.AvgSpeed = sessionMsg.AvgSpeed
			data.AvgCadence = sessionMsg.AvgCadence * 2
			data.TotalCalories = sessionMsg.TotalCalories
			foundSession = true

		case mesgnum.Record:
			// Сообщение Record - детальные данные каждой секунды
			//recordMsg := mesgdef.NewRecord(&message)
			record := entity.RecordData{}
			data.Records = append(data.Records, record)
			// // Если нашли оба ключевых сообщения, можно прекратить поиск
			if foundSession && foundActivity {
				break
			}
		}
	}

	// Проверяем что нашли минимально необходимые данные
	if !foundSession {
		return nil, fmt.Errorf("не найдено сообщение Session в FIT-файле")
	}

	return data, nil
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, w entity.Workout) error {
	return s.Activity.CreateWorkout(ctx, w)
}

func (s *WorkoutService) GetWorkout(ctx context.Context, id int64) (*entity.Workout, error) {
	return &entity.Workout{}, nil
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, u dto.UpdateWorkout) error {
	return s.Activity.UpdateWorkout(ctx, u)
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, id int64) error {
	// TODO: Реализовать удаление из базы данных
	return nil
}
