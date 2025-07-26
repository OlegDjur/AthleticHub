package main

import (
	"log"
	"net/http"
	"workout/internal/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()

	// Логирование и восстановление после паники
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Logger)

	// Твой CORS middleware
	r.Use(corsMiddleware)

	// Статические файлы
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", handler.HomeHandler)
	r.Get("/api/v1/workouts", handler.WorkoutsHandler)
	r.Get("/api/v1/workouts/{id}", handler.WorkoutHandler)
	r.Post("/api/v1/workouts/upload", handler.UploadHandler)
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
	log.Println("Сервер запущен на порту 8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
