package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Тестовый роут
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Mood Tracker API работает!",
		})
	})

	// Простые роуты для начала
	router.POST("/auth/register", registerHandler)
	router.POST("/auth/login", loginHandler)
	router.POST("/entries", addEntryHandler)

	log.Println("Сервер запущен на порту 8080")
	router.Run(":8080")
}

// Простые обработчики
func registerHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Регистрация - заглушка",
	})
}

func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Вход - заглушка",
	})
}

func addEntryHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Добавление записи - заглушка",
	})
}
