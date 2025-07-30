package server

import (
	"errors"
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/decorator"
	custErrors "github.com/Tihmmm/mr-decorator-core/errors"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/validator"
	"github.com/Tihmmm/mr-decorator/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"log"
	"net/http"
)

type Server interface {
	Start(port string) error
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
	if cfg.ApiKey != "" {
		apiKeyHash, err := pkg.GetArgonHash(cfg.ApiKey, nil)
		if err != nil {
			log.Fatalf("Error getting argon2 hash: %v\n", err)
		}
		cfg.ApiKey = apiKeyHash
	}

	server := &EchoServer{
		cfg: cfg,
		e:   echo.New(),
		v:   v,
		d:   d,
	}
	server.registerRoutes()
	return server
}

func (s *EchoServer) Start(port string) error {
	if s.cfg.RateLimit > 0 {
		s.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(s.cfg.RateLimit))))
	}
	log.Printf("Registered parsers: %s", parser.List())
	log.Printf("Starting server on port %s", port)
	if err := s.e.Start(":" + port); err != nil && !errors.Is(http.ErrServerClosed, err) {
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
	return ctx.String(http.StatusOK, "I am alive!")
}

func (s *EchoServer) DecorateMergeRequest(ctx echo.Context) error {
	apiKey := ctx.Request().Header.Get("Api-Key")
	if s.cfg.ApiKey != "" && !pkg.CheckArgonHash(apiKey, s.cfg.ApiKey) {
		return ctx.String(http.StatusUnauthorized, "API Key is invalid")
	}

	mr := new(models.MRRequest)
	err := ctx.Bind(&mr)
	if err != nil {
		log.Printf("Error binding request: %s", err)
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	if !s.v.IsValidAll(mr) {
		return ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	prsr, err := parser.Get(mr.ArtifactFormat)
	if err != nil {
		if errors.Is(err, &custErrors.FormatError{}) {
			return ctx.String(http.StatusBadRequest, fmt.Sprintf("Parser for format `%s` is not supported or registered", mr.ArtifactFormat))
		} else {
			log.Printf("Error getting parser for format %s: %s", mr.ArtifactFormat, err)
			return ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}

	go s.d.Decorate(mr, prsr)

	return ctx.String(http.StatusAccepted, http.StatusText(http.StatusAccepted))
}
