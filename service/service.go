package service

import (
	"time"

	"github.com/Ralphiko/go_final_project/model"
	"github.com/Ralphiko/go_final_project/repository"
)

type Service struct {
	Repository repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s *Service) CreateTask(task model.Task) (int64, error) {
	return s.Repository.CreateTask(task)
}

func (s *Service) GetTaskByID(id string) (model.Task, error) {
	return s.Repository.GetTaskById(id)
}

func (s *Service) SearchTasks(search string) ([]model.Task, error) {
	return s.Repository.GetTasks(search)
}

func (s *Service) UpdateTask(task model.Task) error {
	return s.Repository.UpdateTask(task)
}

func (s *Service) DeleteTask(id string) error {
	return s.Repository.DeleteTask(id)
}

func (s *Service) MarkTaskAsDone(id string) error {
	return s.Repository.TaskDone(id)
}

func (s *Service) GetNextDate(now time.Time, date string, repeat string) (string, error) {
	return s.NextDate(now, date, repeat)
}
