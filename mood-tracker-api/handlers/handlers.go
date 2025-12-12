package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"mood-tracker/models"
	"mood-tracker/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	App *models.App
}

func NewHandler(app *models.App) *Handler {
	return &Handler{App: app}
}

// Регистрация пользователя
func (h *Handler) Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	username := request.Username
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя не может быть пустым"})
		return
	}

	if _, exists := h.App.Users[username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}

	user := &models.User{
		Username: username,
		Entries:  []models.MoodEntry{},
	}

	h.App.Users[username] = *user
	h.App.CurrentUser = user

	if err := storage.Save(h.App); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Аккаунт успешно создан",
		"username": username,
	})
}

// Вход пользователя
func (h *Handler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	username := request.Username
	if user, exists := h.App.Users[username]; exists {
		h.App.CurrentUser = &user
		c.JSON(http.StatusOK, gin.H{
			"message":  "Вход выполнен успешно",
			"username": username,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
	}
}

// Добавление записи настроения
func (h *Handler) AddEntry(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	var request struct {
		Day   int    `json:"day" binding:"required,min=1,max=31"`
		Month int    `json:"month" binding:"required,min=1,max=12"`
		Mood  string `json:"mood" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if !storage.IsValidDay(request.Day, request.Month) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный день для выбранного месяца"})
		return
	}

	if !storage.IsValidMood(request.Mood) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное настроение"})
		return
	}

	entries := h.App.CurrentUser.Entries
	found := false
	for i := range entries {
		if entries[i].Month == request.Month && entries[i].Day == request.Day {
			entries[i].Mood = request.Mood
			found = true
			break
		}
	}
	if !found {
		entries = append(entries, models.MoodEntry{
			Day:   request.Day,
			Month: request.Month,
			Mood:  request.Mood,
		})
	}

	h.App.CurrentUser.Entries = entries
	storage.UpdateUser(h.App, h.App.CurrentUser)

	if err := storage.Save(h.App); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения"})
		return
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"message": "Запись обновлена"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "Запись добавлена"})
	}
}

// Просмотр записи по дате
func (h *Handler) ViewEntry(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	monthStr := c.Query("month")
	dayStr := c.Query("day")

	if monthStr == "" || dayStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать месяц и день"})
		return
	}

	month, err1 := strconv.Atoi(monthStr)
	day, err2 := strconv.Atoi(dayStr)

	if err1 != nil || err2 != nil || !storage.IsValidMonth(month) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные параметры даты"})
		return
	}

	if !storage.IsValidDay(day, month) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный день для выбранного месяца"})
		return
	}

	for _, entry := range h.App.CurrentUser.Entries {
		if entry.Month == month && entry.Day == day {
			c.JSON(http.StatusOK, gin.H{
				"day":   entry.Day,
				"month": entry.Month,
				"mood":  entry.Mood,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Запись не найдена"})
}

// Общий отчёт за месяц
func (h *Handler) GeneralReport(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	monthStr := c.Query("month")
	if monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать месяц"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || !storage.IsValidMonth(month) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный месяц"})
		return
	}

	counts := make(map[string]int)
	for _, entry := range h.App.CurrentUser.Entries {
		if entry.Month == month {
			counts[entry.Mood]++
		}
	}

	if len(counts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Нет записей за указанный месяц",
			"month":   month,
		})
		return
	}

	// Сортировка настроений
	moods := make([]string, 0, len(counts))
	for mood := range counts {
		moods = append(moods, mood)
	}
	sort.Strings(moods)

	result := make(map[string]int)
	for _, mood := range moods {
		result[mood] = counts[mood]
	}

	c.JSON(http.StatusOK, gin.H{
		"month":  month,
		"counts": result,
	})
}

// Отчёт по конкретному настроению
func (h *Handler) MoodReport(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	monthStr := c.Query("month")
	mood := c.Query("mood")

	if monthStr == "" || mood == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать месяц и настроение"})
		return
	}

	if !storage.IsValidMood(mood) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное настроение"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || !storage.IsValidMonth(month) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный месяц"})
		return
	}

	count := 0
	for _, entry := range h.App.CurrentUser.Entries {
		if entry.Month == month && entry.Mood == mood {
			count++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"month": month,
		"mood":  mood,
		"count": count,
	})
}

// Получение всех записей пользователя
func (h *Handler) GetAllEntries(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": h.App.CurrentUser.Username,
		"entries":  h.App.CurrentUser.Entries,
	})
}

// Смена пользователя
func (h *Handler) SwitchUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if user, exists := h.App.Users[request.Username]; exists {
		h.App.CurrentUser = &user
		c.JSON(http.StatusOK, gin.H{
			"message":  "Пользователь изменен",
			"username": request.Username,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
	}
}

// Получение информации о текущем пользователе
func (h *Handler) GetCurrentUser(c *gin.Context) {
	if h.App.CurrentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Нет активного пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      h.App.CurrentUser.Username,
		"entries_count": len(h.App.CurrentUser.Entries),
	})
}

// Тестовый endpoint
func (h *Handler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API работает!",
		"moods":   models.ValidMoods,
	})
}
