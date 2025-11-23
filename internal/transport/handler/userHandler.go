package handler

import (
	"delimed/internal/transport/dto/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

// handler/tarifflist.go

func (h *Handler) GetTariffsList(c echo.Context) error {
	var req request.CDEKRequestList

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input format",
		})
	}

	ctx := c.Request().Context()

	resp, err := h.service.User().GetTariffsList(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTariffs(c echo.Context) error {
	var req request.CDEKRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input format",
		})
	}

	ctx := c.Request().Context()

	resp, err := h.service.User().GetTarifs(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUserProfileHandler(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	user, err := h.service.User().GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUserProfileHandler(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	err := h.service.User().DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "User has deleted"})
}
