package storage

import (
	"encoding/json"
	"mood-tracker/models"
	"os"
)

const dataFile = "mood_diary.json"

// Загрузка данных
func Load(app *models.App) error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			app.Users = make(map[string]models.User)
			return nil
		}
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&app.Users); err != nil {
		return err
	}
	return nil
}

// Сохранение данных
func Save(app *models.App) error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(app.Users)
}

// Получить количество дней в месяце
func DaysInMonth(month int) int {
	switch month {
	case 2:
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

// Проверка валидности месяца
func IsValidMonth(month int) bool {
	return month >= 1 && month <= 12
}

// Проверка валидности дня для месяца
func IsValidDay(day, month int) bool {
	maxDay := DaysInMonth(month)
	return day >= 1 && day <= maxDay
}

// Проверка валидности настроения
func IsValidMood(mood string) bool {
	for _, validMood := range models.ValidMoods {
		if mood == validMood {
			return true
		}
	}
	return false
}

// Получить пользователя
func GetUser(app *models.App, username string) (*models.User, bool) {
	user, exists := app.Users[username]
	return &user, exists
}

// Добавить/обновить пользователя
func UpdateUser(app *models.App, user *models.User) {
	app.Users[user.Username] = *user
}
