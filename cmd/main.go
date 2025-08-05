package main

import (
	"fmt"
	"log"
	"net/http"
	"workout/internal/config"
	handler "workout/internal/controller"
	"workout/internal/service/activity"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// - Заменить handler WorkoutsHandler на новые handlers
// - Полключтить Postgres
// - Настроить сохранение тренировок в бд
// - Сделать эндпоинт для редактирования тренировки
// - Сделать дефолтное имя тренировки
// - Сделать график с объемом за неделю, месяц, год
// - Вывксти рекорды по трассам
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

	svc := activity.NewWorkoutService()
	h := handler.NewWorkoutHandler(svc)

	// r.Get("/", handler.HomeHandler)
	// r.Get("/api/v1/workouts", handler.WorkoutsHandler)
	// r.Get("/api/v1/workouts/{id}", handler.WorkoutHandler)
	e.POST("/api/v1/workouts/upload", h.UploadHandler)
	// e.POST("/api/v1/workouts", h.CreateWorkout)
	// r.Get("/api/v1/workouts/{id}/pacechart", handler.PaceChartHandler) // Получаем пейс для построения графика темпа

	// Модуль Activity
	// /activities - загрузка новой тренировки
	// /activities/{activityID} - Получение деталей тренировки
	// /users/{userID}/activities - Список тренировок пользователя
	// [PUT] /activities/{activityID} - Обновление тренировки (например, добавление заметок)
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
	log.Printf("Сервер запущен на порту %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, e); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
