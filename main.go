package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/malinatrash/golang-initializr/project_templates"
	"github.com/malinatrash/golang-initializr/templates"
)

// Хранилище сгенерированных проектов
var projectFiles = make(map[string]map[string]string)

type ProjectRequest struct {
	Name         string   `json:"name" form:"name"`
	Dependencies []string `json:"dependencies" form:"dependencies"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", handleIndex)
	e.POST("/generate", handleGenerate)
	e.GET("/download", handleDownload)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}

func handleIndex(c echo.Context) error {
	return templates.Index().Render(c.Request().Context(), c.Response().Writer)
}

func handleGenerate(c echo.Context) error {
	req := new(ProjectRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	// Validate project name
	if req.Name == "" {
		return c.String(http.StatusBadRequest, "Project name is required")
	}

	// Generate project in memory
	files := generateProject(req.Name, req.Dependencies)

	// Create a zip archive with project files
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add files to the archive
	for name, content := range files {
		f, err := zipWriter.Create(name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		return err
	}

	// Set headers
	c.Response().Header().Set(echo.HeaderContentType, "application/zip")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=golang-project.zip")

	// Send the zip file as response
	if _, err := c.Response().Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func handleDownload(c echo.Context) error {
	// Получаем последний сгенерированный проект
	session := c.QueryParam("session")
	if session == "" {
		// Если сессия не указана, берем последний сгенерированный проект
		for s := range projectFiles {
			session = s
			break
		}
	}

	files, ok := projectFiles[session]
	if !ok || len(files) == 0 {
		return c.String(http.StatusNotFound, "Project not found")
	}

	// Создаем zip архив с файлами проекта
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Добавляем файлы в архив
	for name, content := range files {
		f, err := zipWriter.Create(name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	// Закрываем zip writer
	if err := zipWriter.Close(); err != nil {
		return err
	}

	// Устанавливаем заголовки
	c.Response().Header().Set(echo.HeaderContentType, "application/zip")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=golang-project.zip")

	// Отправляем zip файл в ответ
	if _, err := c.Response().Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// generateProject генерирует структуру проекта на основе выбранных зависимостей
func generateProject(name string, dependencies []string) map[string]string {
	// Создаем конфигурацию проекта
	config := &project_templates.ProjectConfig{
		Name:         name,
		Dependencies: dependencies,
	}

	// Генерируем файлы проекта
	files := config.GenerateProject()

	fmt.Printf("Generating project %s with dependencies: %s\n", name, strings.Join(dependencies, ", "))
	fmt.Printf("Generated %d files\n", len(files))

	return files
}
