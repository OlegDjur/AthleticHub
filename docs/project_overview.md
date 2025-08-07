# Социальная сеть для атлетов

## Краткий план проекта

**Концепция:**  
Платформа для трекинга тренировок, здоровья, питания и социального взаимодействия спортсменов (бег, плавание, велосипед, силовые тренировки).

---

## Технологии

- **Язык программирования:** Go (Golang)
- **Фреймворк:** Echo (для HTTP API)
- **База данных:** PostgreSQL
- **Кэширование:** Redis
- **Хранение файлов:** S3 (или совместимые решения)
- **JWT:** для аутентификации и авторизации
- **CORS:** собственный middleware или встроенные решения


---

## Модули и ключевые фичи

- **Пользователи:** регистрация, профили (спортсмен/тренер), аутентификация
- **Трекинг активности:** загрузка .fit/.gpx/.tcx, анализ тренировок (темп, пульс, мощность)
- **Единый дашборд:** метрики (тренировки, КБЖУ, сон, вода, пульс, вес)
- **Социальная лента и блог:** посты, фото, комментарии, лайки
- **Взаимодействие тренер-спортсмен:** планы тренировок, отчёты, комментарии тренера
- **Центр здоровья:** дневник питания, трекер БАДов, медицинские документы
- **Сообщество и события:** друзья, группы, события (тренировочные сборы)
- **Карты и маршруты:** создание/обмен маршрутами, GPS-данные
- **Отчёты:** ежедневные, недельные, месячные, годовые отчёты с метриками

---

## Эндпоинты (API /api/v1)

### Пользователи
- POST /users — Регистрация.
- GET /users/{userID} — Профиль.
- PUT /users/{userID} — Обновление профиля.
- DELETE /users/{userID} — Удаление.
- POST /users/login — Аутентификация.
- POST /users/{userID}/athlete-profile — Профиль спортсмена.
- GET /users/{userID}/athlete-profile — Получение профиля спортсмена.
- POST /users/{userID}/coach-profile — Профиль тренера.
- GET /users/{userID}/coach-profile — Получение профиля тренера.

### Трекинг активности
- POST /activities — Загрузка тренировки.
- GET /activities/{activityID} — Детали тренировки.
- GET /users/{userID}/activities — Список тренировок.
- PUT /activities/{activityID} — Обновление.
- DELETE /activities/{activityID} — Удаление.
- GET /activities/summary — Статистика.
- GET /activities/{activityID}/analysis — Анализ тренировки.
- GET /activities/compare — Сравнение тренировок.
- GET /activities/leaderboard — Лидерборд.
- GET /activities/export — Экспорт в CSV.

### Единый дашборд
- GET /users/{userID}/dashboard — Дашборд.
- GET /users/{userID}/dashboard/summary — Сводка.

### Социальная лента и блог
- POST /posts — Создание поста.
- GET /posts/{postID} — Получение поста.
- GET /users/{userID}/posts — Лента пользователя.
- GET /posts/feed — Общая лента.
- PUT /posts/{postID} — Обновление поста.
- DELETE /posts/{postID} — Удаление.
- POST /posts/{postID}/comments — Комментарий.
- GET /posts/{postID}/comments — Список комментариев.
- DELETE /comments/{commentID} — Удаление комментария.
- POST /posts/{postID}/likes — Лайк.
- DELETE /posts/{postID}/likes — Убрать лайк.

### Взаимодействие Тренер-Спортсмен
- POST /coach-athlete-relations — Запрос на тренировку.
- PUT /coach-athlete-relations/{relationID} — Подтверждение/отклонение.
- GET /coaches/{coachID}/athletes — Спортсмены тренера.
- GET /athletes/{athleteID}/coaches — Тренеры спортсмена.
- POST /training-plans — Создание плана.
- GET /training-plans/{planID} — Детали плана.
- GET /athletes/{athleteID}/training-plans — Планы спортсмена.
- POST /trainings — Создание тренировки.
- GET /trainings/{trainingID} — Детали тренировки.
- GET /athletes/{athleteID}/reports — Отчеты спортсмена.
- POST /reports — Создание отчета.

### Центр здоровья
- POST /nutrition-diaries — Запись питания.
- GET /users/{userID}/nutrition-diaries — Дневник питания.
- PUT /nutrition-diaries/{entryID} — Обновление записи.
- DELETE /nutrition-diaries/{entryID} — Удаление записи.
- POST /supplement-trackers — Запись БАДов.
- GET /users/{userID}/supplement-trackers — Записи БАДов.
- DELETE /supplement-trackers/{trackerID} — Удаление записи.
- POST /medical-documents — Загрузка документа.
- GET /users/{userID}/medical-documents — Список документов.
- GET /medical-documents/{documentID} — Получение документа.
- DELETE /medical-documents/{documentID} — Удаление документа.

### Сообщество и события
- POST /friends — Запрос в друзья.
- PUT /friends/{friendID} — Подтверждение/отклонение.
- GET /users/{userID}/friends — Список друзей.
- DELETE /friends/{friendID} — Удаление друга.
- POST /groups — Создание группы.
- GET /groups/{groupID} — Информация о группе.
- GET /users/{userID}/groups — Список групп.
- POST /groups/{groupID}/members — Добавление участника.
- DELETE /groups/{groupID}/members/{userID} — Удаление участника.
- POST /events — Создание события.
- GET /events/{eventID} — Детали события.
- GET /users/{userID}/events — Список событий.
- POST /events/{eventID}/participants — Регистрация на событие.

### Карты и маршруты
- POST /routes — Создание маршрута.
- GET /routes/{routeID} — Детали маршрута.
- GET /users/{userID}/routes — Список маршрутов.
- PUT /routes/{routeID} — Обновление маршрута.
- DELETE /routes/{routeID} — Удаление маршрута.
- POST /routes/{routeID}/share — Поделиться маршрутом.

### Отчеты
- POST /reports — Создание отчета.
- GET /reports/{reportID} — Детали отчета.
- GET /users/{userID}/reports — Список отчетов.
- PUT /reports/{reportID} — Обновление отчета.
- DELETE /reports/{reportID} — Удаление отчета.
- POST /reports/{reportID}/share — Отправка отчета.
- GET /reports/{reportID}/export — Экспорт в PDF/CSV.
- GET /reports/{reportID}/charts — Данные для графиков.
- GET /reports/{reportID}/insights — ИИ-рекомендации.

---

- ИИ-рекомендации по тренировкам и питанию
- Геймификация (значки, челленджи, лидерборды)
- Интеграция с носимыми устройствами (Garmin, Apple Watch)
- Автоматическое распознавание маршрутов из GPS
- Виртуальный тренер на базе ИИ
