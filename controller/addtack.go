package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ralphiko/go_final_project/model"
	database "github.com/Ralphiko/go_final_project/schared/db"
	"github.com/Ralphiko/go_final_project/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) CreateTask(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Создание задачи: %+v", task)

	id, err := h.service.CreateTask(task)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) GetTaskByID(c *gin.Context) {
	id := c.Query("id")
	log.Printf("Запрос задачи с id: %s", id)

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) GetTasks(c *gin.Context) {
	search := c.Query("search")
	log.Printf("Поиск задач по запросу: %s", search)

	tasks, err := h.service.SearchTasks(search)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func TasksReadGET(c *gin.Context) {
	search := c.Query("search")

	var tasks []model.Task
	var err error

	if len(search) > 0 {
		date, err := time.Parse("20060102", search)
		if err != nil {
			tasks, err = database.SearchTasks(search)
		} else {
			tasks, err = database.SearchTasksByDate(date.Format(model.DatePattern))
		}
	} else {
		tasks, err = database.ReadTasks()
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"tasks": []model.Task{}, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func TaskReadGET(c *gin.Context) {
	taskID := c.Query("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	// Преобразование taskID в целочисленный тип
	id, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := database.ReadTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func TaskUpdatePUT(c *gin.Context) {
	var task model.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON decoding error"})
		return
	}

	if task.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if _, err := time.Parse(model.DatePattern, task.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}

	if len(task.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
		return
	}

	if len(task.Repeat) > 0 {
		if _, err := s.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repeat format"})
			return
		}
	}

	updatedTask, err := database.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}
func TaskDonePOST(c *gin.Context) {
	taskID := c.Query("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	// Преобразование taskID в целочисленный тип
	taskIDInt, err := strconv.Atoi(taskID) // Переименовали переменную id в taskIDInt
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := database.ReadTask(taskIDInt) // Заменили id на taskIDInt
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}

	if len(task.Repeat) == 0 {
		if err := database.DeleteTaskDb(taskID); err != nil { // Изменили task.Id на taskIDInt
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
			return
		}
	} else {
		task.Date, err = s.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get next date"})
			return
		}

		if _, err := database.UpdateTask(task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func TaskDELETE(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	if err := database.DeleteTaskDb(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
