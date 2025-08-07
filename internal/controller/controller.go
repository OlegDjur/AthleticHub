package controller

import (
	"io"
	"net/http"
	"workout/internal/dto"
	"workout/internal/entity"
	"workout/internal/service/activity"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	workoutService *activity.WorkoutService
}

func NewController(workoutService *activity.WorkoutService) *Handler {
	return &Handler{
		workoutService: workoutService,
	}
}

// UploadHandler обрабатывает загрузку FIT-файлов через Echo фреймворк
// с использованием новой библиотеки muktihari/fit
func (h *Handler) UploadHandler(c echo.Context) error {
	// Получаем HTTP запрос из контекста Echo
	r := c.Request()

	// Ограничиваем размер тела запроса до 10MB (10 << 20 = 10 * 2^20 = 10485760 байт)
	// Это защищает сервер от загрузки слишком больших файлов
	r.Body = http.MaxBytesReader(c.Response().Writer, r.Body, 10<<20)

	// Парсим multipart/form-data с ограничением в 10MB
	// Эта функция разбирает форму с файлами на составные части
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		// Если парсинг не удался (файл слишком большой или неверный формат)
		// Возвращаем ошибку 400 "Bad Request" через Echo
		return c.String(http.StatusBadRequest, "Файл слишком большой или неверный формат")
	}

	// Извлекаем файл из поля формы с именем "file"
	// В Echo также можно использовать c.FormFile("file") как альтернативу
	// file - интерфейс для чтения файла
	file, _, err := r.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Ошибка при получении файла")
	}
	defer file.Close()

	// Читаем всё содержимое файла в память
	// io.ReadAll читает все данные из Reader до EOF
	data, err := io.ReadAll(file)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при чтении файла")
	}

	workout, err := h.workoutService.UploadFile(c.Request().Context(), data)
	if err != nil {
		return err
	}

	// Отправляем JSON ответ с созданной тренировкой через Echo
	// c.JSON автоматически устанавливает Content-Type: application/json
	// и HTTP статус 201 "Created" указывает что ресурс успешно создан
	return c.JSON(http.StatusCreated, workout)
}
func (h *Handler) CreateWorkout(c echo.Context) error {
	var workout entity.Workout
	if err := c.Bind(&workout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный формат данных"})
	}

	if err := h.workoutService.CreateWorkout(c.Request().Context(), workout); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, workout)
}

func (h *Handler) UpdateWorkout(c echo.Context) error {
	var request dto.UpdateWorkout

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный формат данных"})
	}

	if err := h.workoutService.UpdateWorkout(c.Request().Context(), request); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")
}
