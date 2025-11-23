package httpserver

import (
	"delimed/internal/transport/handler"
	customMiddleware "delimed/internal/transport/middleware"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewServer(jwtSecret []byte, h *handler.Handler, addr string, writeTimeout, readTimeout, idleTimeout time.Duration) *http.Server {
	router := echo.New()

	router.Use(middleware.BodyLimit("10M"))

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
		AllowCredentials: true,
	}))

	// Swagger UI
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.GET("/tariffslist", h.GetTariffsList)
	router.GET("/tariffs", h.GetTariffs)
	router.POST("/delivery/calculate", h.CalculateDeliveryOptions)
	router.POST("/register", h.RegisterUserHandler)
	router.POST("/login", h.LoginUserHandler)

	protected := router.Group("/api")
	{
		protected.Use(customMiddleware.AuthRequired(jwtSecret))

		user := protected.Group("/user")
		{
			user.GET("", h.GetUserProfileHandler)
			user.DELETE("", h.DeleteUserProfileHandler)
		}
	}

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}
}
