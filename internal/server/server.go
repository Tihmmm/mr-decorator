package server

import (
	"errors"
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/decorator"
	"github.com/Tihmmm/mr-decorator/internal/models"
	"github.com/Tihmmm/mr-decorator/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

type Server interface {
	Start() error
	Liveliness(ctx echo.Context) error
	DecorateMergeRequest(ctx echo.Context) error
}

type EchoServer struct {
	cfg config.ServerConfig
	e   *echo.Echo
	v   validator.Validator
	d   decorator.Decorator
}

func NewEchoServer(cfg config.ServerConfig, v validator.Validator, d decorator.Decorator) Server {
	server := &EchoServer{
		cfg: cfg,
		e:   echo.New(),
		v:   v,
		d:   d,
	}
	server.registerRoutes()
	return server
}

func (s *EchoServer) Start() error {
	s.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(3)))
	if err := s.e.Start(s.cfg.Port); err != nil && !errors.Is(http.ErrServerClosed, err) {
		log.Fatalf("Server shutdown occured: %s", err)
		return err
	}
	return nil
}

func (s *EchoServer) registerRoutes() {
	s.e.GET("/liveliness", s.Liveliness)
	s.e.POST("/decorate-merge-request", s.DecorateMergeRequest)
}

func (s *EchoServer) Liveliness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{Status: http.StatusText(http.StatusOK)})
}

func (s *EchoServer) DecorateMergeRequest(ctx echo.Context) error {
	mrRequest := new(models.MRRequest)
	err := ctx.Bind(&mrRequest)
	if err != nil {
		log.Printf("Error binding request: %s", err)
		return ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	if !s.v.IsValidAll(mrRequest) {
		return ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	go s.d.DecorateServer(mrRequest)

	return ctx.JSON(http.StatusAccepted, http.StatusText(http.StatusAccepted))
}
