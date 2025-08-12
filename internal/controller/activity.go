package controller

import (
	"github.com/gofrs/uuid/v5"
	"io"
	"net/http"
	"workout/internal/controller/mapper"
	"workout/internal/dto"
	"workout/internal/service/activity"
	"workout/internal/service/auth"

	"github.com/labstack/echo/v4"
)

const ctxKeyClaims = "claims"
const ctxKeySub = "sub"

func CurrentUser(c echo.Context) (*UserClaims, bool) {
	v := c.Get(ctxKeyClaims)
	if v == nil {
		return nil, false
	}
	uc, ok := v.(*UserClaims)
	return uc, ok
}

type Handler struct {
	workoutService *activity.WorkoutService
	auth           *auth.AuthService
}

func NewController(workoutService *activity.WorkoutService, authService *auth.AuthService) *Handler {
	return &Handler{
		workoutService: workoutService,
		auth:           authService,
	}
}

func (h *Handler) GetWorkouts(e echo.Context) error {
	user, ok := CurrentUser(e)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no auth context")
	}

	workout, err := h.workoutService.GetWorkouts(e.Request().Context(), user.UID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	result := []*dto.WorkoutDTO{}

	for _, w := range workout {
		result = append(result, mapper.ConvertWorkoutToDTO(w))
	}

	return e.JSON(http.StatusCreated, result)
}

// UploadHandler обрабатывает загрузку FIT-файлов через Echo фреймворк
// с использованием новой библиотеки muktihari/fit
func (h *Handler) UploadHandler(e echo.Context) error {
	// Получаем HTTP запрос из контекста Echo
	r := e.Request()

	// Ограничиваем размер тела запроса до 10MB (10 << 20 = 10 * 2^20 = 10485760 байт)
	// Это защищает сервер от загрузки слишком больших файлов
	r.Body = http.MaxBytesReader(e.Response().Writer, r.Body, 10<<20)

	// Парсим multipart/form-data с ограничением в 10MB
	// Эта функция разбирает форму с файлами на составные части
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		// Если парсинг не удался (файл слишком большой или неверный формат)
		// Возвращаем ошибку 400 "Bad Request" через Echo
		return e.String(http.StatusBadRequest, "Файл слишком большой или неверный формат")
	}

	// Извлекаем файл из поля формы с именем "file"
	// В Echo также можно использовать c.FormFile("file") как альтернативу
	// file - интерфейс для чтения файла
	file, _, err := r.FormFile("file")
	if err != nil {
		return e.String(http.StatusBadRequest, "Ошибка при получении файла")
	}
	defer file.Close()

	// Читаем всё содержимое файла в память
	// io.ReadAll читает все данные из Reader до EOF
	data, err := io.ReadAll(file)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Ошибка при чтении файла")
	}

	workout, err := h.workoutService.UploadFile(e.Request().Context(), data)
	if err != nil {
		return err
	}

	// Отправляем JSON ответ с созданной тренировкой через Echo
	// c.JSON автоматически устанавливает Content-Type: application/json
	// и HTTP статус 201 "Created" указывает что ресурс успешно создан
	return e.JSON(http.StatusCreated, workout)
}
func (h *Handler) CreateWorkout(e echo.Context) error {
	user, ok := CurrentUser(e)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no auth context")
	}

	var workout dto.WorkoutDTO
	if err := e.Bind(&workout); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный формат данных"})
	}

	u, err := uuid.FromString(user.UID)
	if err != nil {
		return err
	}
	workout.UserID = u
	response, err := h.workoutService.CreateWorkout(e.Request().Context(), workout)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusCreated, response)
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
