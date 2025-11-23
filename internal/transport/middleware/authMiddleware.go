package middleware

import (
	"delimed/internal/utils/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthRequired(jwtSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "missing bearer"})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ValidateToken(tokenStr, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid token"})
			}

			uid := claims.ID

			c.Set("user_id", uid.String())
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role == "" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "no role"})
			}

			if _, allowed := allowed[role]; !allowed {
				return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "forbidden"})
			}

			return next(c)
		}
	}
}
