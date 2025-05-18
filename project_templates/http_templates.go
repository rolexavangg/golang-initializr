package project_templates

const httpHandlerTemplate = `package http

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	
	"github.com/labstack/echo/v4"
	"{{.Name}}/internal/config"
)


var Module = fx.Options(
	fx.Provide(NewUserHandler),
)


func RegisterRoutes(server *echo.Echo, userHandler *UserHandler) {
	api := server.Group("/api")
	
	
	users := api.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("/:id", userHandler.GetByID)
	users.GET("", userHandler.List)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)
}
`

const httpServerTemplate = `package http

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type Server struct {
	echo   *echo.Echo
	logger *zap.Logger
}


func NewServer(e *echo.Echo, logger *zap.Logger) *Server {
	return &Server{
		echo:   e,
		logger: logger,
	}
}
`

const httpUserHandlerTemplate = `package http

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/domain"
)


type UserHandler struct {
	useCase domain.UserUseCase
	logger  *zap.Logger
}


func NewUserHandler(useCase domain.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		logger:  logger,
	}
}


func (h *UserHandler) Create(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if err := h.useCase.Create(user); err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	
	return c.JSON(http.StatusCreated, user)
}


func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	
	user, err := h.useCase.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user"})
	}
	
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	return c.JSON(http.StatusOK, user)
}


func (h *UserHandler) List(c echo.Context) error {
	users, err := h.useCase.List()
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list users"})
	}
	
	return c.JSON(http.StatusOK, users)
}


func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	user.ID = id
	
	if err := h.useCase.Update(user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	
	return c.JSON(http.StatusOK, user)
}


func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	
	if err := h.useCase.Delete(id); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	
	return c.NoContent(http.StatusNoContent)
}
`
