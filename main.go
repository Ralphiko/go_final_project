package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ralphiko/go_final_project/controller"
	"github.com/Ralphiko/go_final_project/service"
	"github.com/gin-gonic/gin"
)

var (
	webDir = "./web"
	router = gin.Default()
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func main() {
	portFromEnv := os.Getenv("TODO_PORT")
	var port int
	if portFromEnv != "" {
		p, err := strconv.Atoi(portFromEnv)
		if err != nil {
			log.Fatalf("Failed to parse TODO_PORT: %v", err)
		}
		port = p
	} else {
		port = 7540
	}

	service := service.NewService(db)
	handler := NewHandler(service)

	api := router.Group("/api")
	{
		api.GET("/tasks", controller.TasksReadGET)
		api.POST("/task", controller.CreateTask)
		api.GET("/task", controller.TaskReadGET)
		api.PUT("/task", controller.TaskUpdatePUT)
		api.DELETE("/task", controller.TaskDELETE)
		api.POST("/task/done", controller.TaskDonePOST)
	}

	// Регистрация маршрутов для API
	router.GET("api/nextdate", service.NextDateHandler)

	// Простой файл-сервер для отдачи статических файлов
	router.NoRoute(func(c *gin.Context) {
		http.FileServer(http.Dir(webDir)).ServeHTTP(c.Writer, c.Request)
	})

	// Запуск сервера
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on port %d\n", port)
	log.Fatal(router.Run(addr))
}
