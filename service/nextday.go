package service

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func addDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

func addYears(date time.Time, years int) time.Time {
	return date.AddDate(years, 0, 0)
}

func getNextDay(now time.Time, date string, repeat string) (string, error) {
	days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
	if err != nil {
		return "", err
	}

	if days > 400 {
		return "", errors.New("number of days exceeds the maximum limit of 400")
	}

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	newDate := addDays(parsedDate, days)

	for newDate.Before(now) {
		newDate = addDays(newDate, days)
	}

	return newDate.Format("20060102"), nil
}

func getNextYear(now time.Time, date string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	newDate := addYears(parsedDate, 1)

	for newDate.Before(now) {
		newDate = addYears(newDate, 1)
	}

	return newDate.Format("20060102"), nil
}

func (s *Service) NextDate(now time.Time, date string, repeat string) (string, error) {
	if strings.HasPrefix(repeat, "d") {
		return getNextDay(now, date, repeat)
	} else if strings.Contains(repeat, "y") {
		return getNextYear(now, date)
	}
	return "", errors.New("repeat wrong format")
}

func (s *Service) NextDateHandler(c *gin.Context) {
	now := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	log.Println(now, date, repeat)

	if now == "" || date == "" || repeat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствуют обязательные параметры"})
		return
	}

	parsedNow, err := time.Parse("20060102", now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный параметр 'now'"})
		return
	}

	nextDate, err := s.NextDate(parsedNow, date, repeat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, nextDate) // Исправлено на возвращение строки без JSON-обёртки
}
