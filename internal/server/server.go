package server

import (
	"errors"
	"fmt"
	"github.com/Tihmmm/mr-decorator/internal/client"
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/models"
	"github.com/Tihmmm/mr-decorator/internal/parser"
	"github.com/Tihmmm/mr-decorator/internal/validator"
	"github.com/Tihmmm/mr-decorator/pkg/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"
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
	c   client.Client
	p   parser.Parser
}

const waitTime = 4 * time.Second

func NewEchoServer(cfg config.ServerConfig, v validator.Validator, c client.Client, p parser.Parser) Server {
	server := &EchoServer{
		cfg: cfg,
		e:   echo.New(),
		v:   v,
		c:   c,
		p:   p,
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
	if !s.v.Validate(mrRequest) {
		return ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	go s.decMR(mrRequest)

	return ctx.JSON(http.StatusAccepted, http.StatusText(http.StatusAccepted))
}

func (s *EchoServer) decMR(mrRequest *models.MRRequest) {
	time.Sleep(waitTime)

	fmt.Printf("%s Started processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.JobId)

	artifactsDir, err := s.c.GetArtifact(mrRequest.ProjectId, mrRequest.JobId, mrRequest.ArtifactFileName, mrRequest.AuthToken)
	if err != nil {
		return
	}
	defer file.DeleteDirectory(artifactsDir)

	note, err := s.p.Parse(mrRequest.ArtifactFormat, mrRequest.ArtifactFileName, artifactsDir, mrRequest.VulnerabilityMgmtId)
	if err != nil {
		return
	}

	err = s.c.SendNote(note, mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.AuthToken)
	if err != nil {
		return
	}

	fmt.Printf("%s Finished processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.JobId)
}
