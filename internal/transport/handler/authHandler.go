package handler

import (
	"delimed/internal/transport/dto/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterUserHandler(c echo.Context) error {

	var req request.SignUpInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
	}

	if err := h.service.Auth().RegisterUser(req); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Registration successful."})
}

func (h *Handler) LoginUserHandler(c echo.Context) error {
	var req request.SignInInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
	}

	token, err := h.service.Auth().LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
