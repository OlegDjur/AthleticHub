package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"workout/internal/adapters/postgres"
	"workout/internal/config"
	handler "workout/internal/controller"
	"workout/internal/service/activity"
	"workout/internal/service/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// [] добавить возможность загрузки нескольких тренировок
// [] вынести в пакет utils convertDistance
// [] Создать структуру ответа после парсинга фит файлаы
// [x] Заменить handler WorkoutsHandler на новые handlers
// [x] Полключтить Postgres
// [x] Настроить сохранение тренировок в бд
// - Сделать эндпоинт для редактирования тренировки
// - Сделать дефолтное имя тренировки
// - Сделать график с объемом за неделю, месяц, год
// - Вывести рекорды по трассам
// - Сделать профиль юзера
// - Создать модуль авторизации
// - Настроить GraceFull shutdown

func corsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusOK)
		}
		return next(c)
	}
}

func main() {
	cfg, err := config.LoadConfig("config/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации: ", err)
	}
	fmt.Printf("%+v", cfg)

	e := echo.New()

	// Логирование и восстановление после паники
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(corsMiddleware)

	repo, err := postgres.NewPostgresAdapter(cfg.Database)
	if err != nil {
		return
	}

	auth := auth.NewAuthService(repo, cfg.Auth.TokenTTL)
	svc := activity.NewWorkoutService(repo)
	h := handler.NewController(svc, auth)

	e.POST("/login", h.Login)
	e.POST("/register", h.Register)

	api := e.Group("/api", handler.JWTMiddleware([]byte("")))

	api.POST("/v1/workout/upload", h.UploadHandler)
	api.POST("/v1/workout", h.CreateWorkout)
	api.PUT("/v1/workout/{id}", h.UpdateWorkout) // Обновление тренировки (например, добавление заметок)
	api.GET("/v1/workouts", h.GetWorkouts)
	// r.Get("/api/v1/workouts/{id}", handler.WorkoutHandler)
	// r.Get("/api/v1/workouts/{id}/pacechart", handler.PaceChartHandler) // Получаем пейс для построения графика темпа
	// r.Get("/", handler.HomeHandler)

	// Модуль Activity
	// /activities - загрузка новой тренировки
	// /activities/{activityID} - Получение деталей тренировки
	// /users/{userID}/activities - Список тренировок пользователя
	// [DELETE] /activities/{activityID} - Удаление тренировки
	// [GET] /activities/summary - Агрегированная статистика тренировок
	// [GET] /activities/{activityID}/analysis - Детальный анализ тренировки
	// [GET] /api/v1/activities/compare?activityIDs=1,2,3 — Сравнение нескольких тренировок.
	// [GET] /api/v1/activities/leaderboard?groupID=123 — Лидерборд по тренировкам в группе.
	// [GET] /api/v1/activities/{activityID}/analysis — Детальный анализ тренировки (сравнение с прошлыми, тренды).

	// Модуль Reports
	// POST /api/v1/reports - Создание отчета (указать тип: день/неделя/месяц/год, метрики)
	// GET /api/v1/users/{userID}/reports — Список отчетов с фильтрами (?type=daily, ?date=2025-07-07)
	// PUT /api/v1/reports/{reportID} — Обновление отчета (например, добавление заметок)
	// GET /api/v1/reports/{reportID} — Детали отчета
	// GET /api/v1/reports/{reportID}/charts — Получение данных для графиков (например, пульс за неделю)
	// GET /api/v1/reports/{reportID}/insights — Получение ИИ-рекомендаций (например, "Сон менее 7 часов в 4 из 7 дней, добавьте отдых")
	// POST /api/v1/reports/{reportID}/share — Отправка отчета тренеру
	// GET /api/v1/reports/{reportID}/achievements — Список достижений в отчете (например, "5 тренировок за неделю")
	// POST /api/v1/reports/achievements — Создание пользовательских достижений

	// GET /api/v1/coaches/{coachID}/reports — Список отчетов от всех спортсменов тренера.

	// Запускаем сервер
	log.Printf("Сервер запущен на порту %d\n", cfg.Server.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Server.Port), e); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
