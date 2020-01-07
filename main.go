package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Task struct {
	Id          uint   `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskHandler struct {
	DB *gorm.DB
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	h := TaskHandler{}
	h.Initialize()

	e.GET("/tasks", h.GetAllTask)
	e.POST("/tasks", h.SaveTask)
	e.GET("/tasks/:id", h.GetTask)
	e.PUT("/tasks/:id", h.UpdateTask)
	e.DELETE("/tasks/:id", h.DeleteTask)

	e.Logger.Fatal(e.Start(":8081"))
}

func (h *TaskHandler) Initialize() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Task{})

	h.DB = db
}

func (h *TaskHandler) GetAllTask(c echo.Context) error {
	tasks := []Task{}

	h.DB.Find(&tasks)

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	id := c.Param("id")
	task := Task{}

	if err := h.DB.Find(&task, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) SaveTask(c echo.Context) error {
	task := Task{}

	if err := c.Bind(&task); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&task).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	id := c.Param("id")
	task := Task{}

	if err := h.DB.Find(&task, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Bind(&task); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&task).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id := c.Param("id")
	task := Task{}

	if err := h.DB.Find(&task, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := h.DB.Delete(&task).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
