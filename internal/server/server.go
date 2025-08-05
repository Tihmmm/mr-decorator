package server

import (
	"errors"
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/decorator"
	custErrors "github.com/Tihmmm/mr-decorator-core/errors"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/validator"
	"github.com/Tihmmm/mr-decorator/config"
	"github.com/Tihmmm/mr-decorator/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"log"
	"net/http"
)

type Server interface {
	Start(port string) error
	HealthCheck(ctx echo.Context) error
	DecorateMergeRequest(ctx echo.Context) error
}

type EchoServer struct {
	cfg config.ServerConfig
	e   *echo.Echo
	v   validator.Validator
	d   decorator.Decorator
}

func NewEchoServer(cfg config.ServerConfig, v validator.Validator, d decorator.Decorator) Server {
	e := echo.New()
	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit))))
	}
	server := &EchoServer{
		cfg: cfg,
		e:   e,
		v:   v,
		d:   d,
	}

	e.GET("/healthcheck", server.HealthCheck)

	groupInternal := e.Group("/internal")
	if cfg.ApiKey != "" {
		var err error
		apiKeyHash, err = pkg.GetArgonHash(cfg.ApiKey, nil)
		if err != nil {
			log.Fatalf("Error getting argon2 hash: %v\n", err)
		}
		cfg.ApiKey = apiKeyHash

		groupInternal.Use(authMiddleware)
	}
	groupInternal.POST("/decorate-merge-request", server.DecorateMergeRequest)

	return server
}

func (s *EchoServer) Start(port string) error {
	log.Printf("Registered parsers: %s", parser.List())
	log.Printf("Starting server on port %s", port)

	return s.e.Start(":" + port)
}

func (s *EchoServer) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "I am alive!")
}

func (s *EchoServer) DecorateMergeRequest(ctx echo.Context) error {
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
