package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
	"workout/internal/models"
	"workout/internal/service"

	"github.com/tormoder/fit"
)

// Параметры как у Garmin
const (
	intervalDistance = 0.001 // Интервал 5м для расчета темпа
	lookbackDistance = 0.1   // Смотреть назад на 100м для усреднения
	maxPace          = 15.0  // Максимальный темп (мин/км)
	minPace          = 2.5   // Минимальный темп (мин/км)
	maxPaceChange    = 2.0   // Максимальное изменение темпа между точками (мин/км)
)

type Handler struct {
	workoutService *service.WorkoutService
}

func NewWorkoutHandler(workoutService *service.WorkoutService) *Handler {
	return &Handler{
		workoutService: workoutService,
	}
}

// Обработчик корневого маршрута
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Привет! Это API для мониторинга тренировок."))
}

func GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Возвращаем список всех тренировок
	// Пока что тестовые данные (заглушка)
	json.NewEncoder(w).Encode(models.Workouts)
}

func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Создаем новую тренировку, сохраняется пока что просто в памяти
	var newWorkout models.Workout
	if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Генерируется id
	newWorkout.ID = fmt.Sprintf("%d", len(models.Workouts)+1)
	models.Workouts = append(models.Workouts, newWorkout)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWorkout)
}

// Обработчик для списка тренировок
func WorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	switch r.Method {
	case "GET":
		// Возвращаем список всех тренировок
		// Пока что тестовые данные (заглушка)
		json.NewEncoder(w).Encode(models.Workouts)
	case "POST":
		// Создаем новую тренировку, сохраняется пока что просто в памяти
		var newWorkout models.Workout
		if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}
		newWorkout.ID = fmt.Sprintf("%d", len(models.Workouts)+1)
		models.Workouts = append(models.Workouts, newWorkout)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newWorkout)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик для операций с конкретной тренировкой
func WorkoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/workouts/")

	switch r.Method {
	case "GET":
		// Получаем тренировку по ID
		for _, workout := range models.Workouts {
			if workout.ID == id {
				json.NewEncoder(w).Encode(workout)
				return
			}
		}
		http.Error(w, "Тренировка не найдена", http.StatusNotFound)
	case "PUT":
		// Обновляем тренировку
		var updatedWorkout models.Workout
		if err := json.NewDecoder(r.Body).Decode(&updatedWorkout); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}
		for i, workout := range models.Workouts {
			if workout.ID == id {
				updatedWorkout.ID = id
				models.Workouts[i] = updatedWorkout
				json.NewEncoder(w).Encode(updatedWorkout)
				return
			}
		}
		http.Error(w, "Тренировка не найдена", http.StatusNotFound)
	case "DELETE":
		// Удаляем тренировку
		for i, workout := range models.Workouts {
			if workout.ID == id {
				models.Workouts = append(models.Workouts[:i], models.Workouts[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Тренировка не найдена", http.StatusNotFound)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик загрузки файлов
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Ограничиваем размер файла (например, 10MB)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Файл слишком большой или неверный формат", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Парсинг FIT-файла
	fitFile, err := fit.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Ошибка при парсинге FIT-файла", http.StatusBadRequest)
		return
	}

	// Извлекаем данные из FIT-файла
	activity, err := fitFile.Activity()
	if err != nil {
		http.Error(w, "Ошибка при извлечении активности", http.StatusInternalServerError)
		return
	}

	pace := calculatePace(activity.Sessions[0].TotalDistance, activity.Sessions[0].TotalTimerTime)

	log.Printf("[uploadHandler] TotalDistance=%d, TotalTimerTime=%d", activity.Sessions[0].TotalDistance, activity.Sessions[0].TotalTimerTime)
	log.Printf("[uploadHandler] Рассчитанный темп: %.3f", pace)

	distance := float64(activity.Sessions[0].TotalDistance) / 100000.0

	// Формируем объект тренировки
	newWorkout := models.Workout{
		ID:        fmt.Sprintf("%d", len(models.Workouts)+1),
		Date:      activity.Sessions[0].Timestamp.Format("2006-01-02"),
		Distance:  fmt.Sprintf("%.2f", distance),                                 // метры в км
		Duration:  secondsToHMS(int(activity.Sessions[0].TotalTimerTime) / 1000), // секунды
		Pace:      formatPace(pace),
		HeartRate: int(activity.Sessions[0].AvgHeartRate),
		// GPS:       ExtractPaceChart(activity.Records),
	}

	// Сохраняем тренировку
	models.Workouts = append(models.Workouts, newWorkout)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWorkout)
}

// Функция для конвертации секунд в формат "часы:минуты:секунды"
func secondsToHMS(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// Вспомогательная функция для вычисления темпа (мин/км)
func calculatePace(distance uint32, duration uint32) float64 {
	if distance == 0 || duration == 0 {
		return 0
	}

	// Дистанция в километры:
	// FIT distance в сантиметрах
	// 100000 см = 1 км
	// → надо делить на 100000
	km := float64(distance) / 100000.0

	// Время в минутах
	minutes := float64(duration) / 1000.0 / 60.0

	return minutes / km
}

func formatPace(pace float64) string {
	if pace == 0 {
		return "-"
	}

	// Проверяем, что темп в разумных пределах
	if pace < 1.0 || pace > 20.0 {
		log.Printf("[formatPace] Недопустимый темп: %.3f", pace)
		return "-"
	}

	mins := int(pace)
	secs := int((pace - float64(mins)) * 60)

	// Проверяем, что секунды не превышают 59
	if secs >= 60 {
		mins += secs / 60
		secs = secs % 60
	}

	// Дополнительная проверка - секунды не могут быть больше 59
	if secs > 59 {
		log.Printf("[formatPace] ОШИБКА: секунды превышают 59: %d", secs)
		secs = 59
	}

	log.Printf("[formatPace] pace=%.3f -> %d:%02d мин/км", pace, mins, secs)
	return fmt.Sprintf("%d:%02d мин/км", mins, secs)
}

// Функция медианного сглаживания
func median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// Функция для расчета среднего пульса
func averageHeartRate(heartRates []int) int {
	if len(heartRates) == 0 {
		return 0
	}
	sum := 0
	validCount := 0
	for _, hr := range heartRates {
		if hr > 0 { // Игнорируем нулевые значения пульса
			sum += hr
			validCount++
		}
	}
	if validCount == 0 {
		return 0
	}
	return sum / validCount
}

// ExtractPaceChart- не используется, необходим для постороения графика темпа
func ExtractPaceChart(records []*fit.RecordMsg) []models.PacePoint {
	var result []models.PacePoint

	log.Printf("[PaceChart] Начинаем обработку %d записей", len(records))

	if len(records) < 20 {
		log.Printf("[PaceChart] Недостаточно записей: %d", len(records))
		return result
	}

	// Находим первую валидную запись
	startIdx := 0
	for i, record := range records {
		if record != nil && record.Distance > 0 {
			startIdx = i
			break
		}
	}

	// Конвертируем все расстояния в км и собираем данные о пульсе
	var distances []float64
	var times []time.Time
	var heartRates []int
	for i := startIdx; i < len(records); i++ {
		if records[i] != nil && records[i].Distance > 0 {
			rawDistance := float64(records[i].Distance)
			convertedDistance := rawDistance / 100000.0
			distances = append(distances, convertedDistance)
			times = append(times, records[i].Timestamp)
			heartRates = append(heartRates, int(records[i].HeartRate))

			// Логируем первые несколько записей для отладки
			if i < startIdx+5 {
				log.Printf("[PaceChart] Запись %d: rawDistance=%.0f, convertedDistance=%.3f км, HR=%d",
					i, rawDistance, convertedDistance, records[i].HeartRate)
			}
		}
	}

	if len(distances) < 2 {
		return result
	}

	var lastValidPace float64
	var paceWindow []float64 // Окно для сглаживания

	// Рассчитываем темп на каждом интервале
	// Логируем последние 5 точек distances
	for i := len(distances) - 5; i < len(distances); i++ {
		if i >= 0 {
			log.Printf("[PaceChart] LAST distance[%d]=%.3f", i, distances[i])
		}
	}

	maxDistance := distances[len(distances)-1]
	// Если известно реальное значение дистанции (например, totalDistance), обрезаем maxDistance
	if len(records) > 0 && records[len(records)-1] != nil && records[len(records)-1].Distance > 0 {
		realMaxDistance := float64(records[len(records)-1].Distance) / 100000.0
		if maxDistance > realMaxDistance {
			log.Printf("[PaceChart] Обрезаем maxDistance: %.3f -> %.3f", maxDistance, realMaxDistance)
			maxDistance = realMaxDistance
		}
	}
	currentDistance := intervalDistance
	for currentDistance <= maxDistance {
		if currentDistance > maxDistance {
			break
		}
		// Для последних точек, если до конца дистанции осталось меньше lookbackDistance, используем startDist = 0
		startDist := currentDistance - lookbackDistance
		if startDist < 0 {
			startDist = 0
		}

		// Находим индексы точек
		var startIdx, endIdx int
		for i, dist := range distances {
			if dist >= startDist {
				startIdx = i
				break
			}
		}
		for i, dist := range distances {
			if dist >= currentDistance {
				endIdx = i
				break
			}
		}

		if endIdx <= startIdx || startIdx >= len(distances)-1 {
			currentDistance += intervalDistance
			continue
		}

		// Рассчитываем общее время и расстояние для интервала
		totalDistance := distances[endIdx] - distances[startIdx]
		totalTime := times[endIdx].Sub(times[startIdx]).Seconds()

		if totalDistance < 0.01 || totalTime <= 0 { // Минимум 10м и положительное время
			currentDistance += intervalDistance
			continue
		}

		// Рассчитываем темп
		pace := totalTime / (60 * totalDistance)

		// Фильтруем аномальные значения
		if pace < minPace || pace > maxPace {
			log.Printf("[PaceChart] Аномальный темп: %.2f мин/км на %.2f км", pace, totalDistance)
			currentDistance += intervalDistance
			continue
		}

		// Фильтруем резкие изменения темпа
		if lastValidPace > 0 && math.Abs(pace-lastValidPace) > maxPaceChange {
			log.Printf("[PaceChart] Слишком резкое изменение темпа: %.2f -> %.2f мин/км", lastValidPace, pace)
			// Используем среднее значение вместо резкого скачка
			pace = (lastValidPace + pace) / 2
		}

		// Добавляем в окно сглаживания
		paceWindow = append(paceWindow, pace)
		if len(paceWindow) > 20 { // Окно из 20 точек
			paceWindow = paceWindow[1:]
		}
		// Медианное сглаживание
		var smoothedPace float64
		if len(paceWindow) >= 5 {
			smoothedPace = median(paceWindow)
		} else {
			smoothedPace = pace
		}

		// Округляем дистанцию и темп
		distanceKm := math.Round(currentDistance*100) / 100

		// Округляем темп до целых секунд (например, 4.93 -> 4.56)
		mins := int(smoothedPace)
		secs := int((smoothedPace - float64(mins)) * 60)
		roundedPace := float64(mins) + float64(secs)/100

		// Рассчитываем средний пульс для интервала
		var intervalHeartRates []int
		for i := startIdx; i <= endIdx && i < len(heartRates); i++ {
			if heartRates[i] > 0 {
				intervalHeartRates = append(intervalHeartRates, heartRates[i])
			}
		}
		avgHeartRate := averageHeartRate(intervalHeartRates)

		result = append(result, models.PacePoint{
			DistanceKm: distanceKm,
			Pace:       roundedPace,
			HeartRate:  avgHeartRate,
			Time:       times[endIdx],
		})

		// Проверяем, что темп в разумных пределах
		if smoothedPace < 2.0 || smoothedPace > 15.0 {
			log.Printf("[PaceChart] ВНИМАНИЕ: Недопустимый темп %.3f на дистанции %.2f км", smoothedPace, distanceKm)
		}

		// Проверяем, что темп не превышает 10 минут на километр (очень медленный бег)
		if smoothedPace > 10.0 {
			log.Printf("[PaceChart] ОШИБКА: Слишком медленный темп %.3f мин/км", smoothedPace)
			continue // Пропускаем эту точку
		}

		lastValidPace = smoothedPace

		currentDistance += intervalDistance
	}

	log.Printf("[PaceChart] Обработано %d точек темпа", len(result))
	return result
}

// func PaceChartHandler(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")

// 	var workout models.Workout

// 	for _, w := range models.Workouts {
// 		if w.ID == id {
// 			workout = w
// 		}
// 	}

// 	// Логируем первые несколько точек для отладки
// 	if len(workout.GPS) > 0 {
// 		log.Printf("[paceChartHandler] Первые 3 точки темпа:")
// 		for i := 0; i < 3 && i < len(workout.GPS); i++ {
// 			point := workout.GPS[i]
// 			log.Printf("  [%d] distance=%.2f km, pace=%.3f (тип: %T)",
// 				i, point.DistanceKm, point.Pace, point.Pace)
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(workout.GPS)
// }
