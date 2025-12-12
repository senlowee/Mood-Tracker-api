package models

// Настроения
const (
	MoodSad     = "Грустное"
	MoodCalm    = "Спокойное"
	MoodNeutral = "Нейтральное"
	MoodHappy   = "Радостное"
	MoodAngry   = "Злое"
)

var ValidMoods = []string{MoodSad, MoodCalm, MoodNeutral, MoodHappy, MoodAngry}

// Месяцы
var Months = []string{
	"январь", "февраль", "март", "апрель", "май", "июнь",
	"июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь",
}

// Запись настроения
type MoodEntry struct {
	Day   int    `json:"day"`
	Month int    `json:"month"`
	Mood  string `json:"mood"`
}

// Пользователь
type User struct {
	Username string      `json:"username"`
	Entries  []MoodEntry `json:"entries"`
}

// Приложение/Хранилище
type App struct {
	Users       map[string]User `json:"users"`
	CurrentUser *User           `json:"-"`
}
