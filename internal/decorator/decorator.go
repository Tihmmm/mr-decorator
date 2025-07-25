package decorator

import (
	"github.com/Tihmmm/mr-decorator/internal/client"
	"github.com/Tihmmm/mr-decorator/internal/models"
	"github.com/Tihmmm/mr-decorator/internal/parser"
	"github.com/Tihmmm/mr-decorator/pkg/file"
	"log"
	"path/filepath"
	"time"
)

type Decorator interface {
	DecorateServer(mrRequest *models.MRRequest) error
	DecorateCli(mrRequest *models.MRRequest) error
}

type MRDecorator struct {
	mode string // either `cli` or `server`
	c    client.Client
	p    parser.Parser
}

func NewDecorator(m string, c client.Client, p parser.Parser) Decorator {
	return &MRDecorator{
		mode: m,
		c:    c,
		p:    p,
	}
}

const waitTime = 4 * time.Second // waiting for artifacts to be loaded

func (d *MRDecorator) DecorateServer(mrRequest *models.MRRequest) error {
	time.Sleep(waitTime)

	log.Printf("%s Started processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.JobId)

	artifactsDir, err := d.c.GetArtifact(mrRequest.ProjectId, mrRequest.JobId, mrRequest.ArtifactFileName, mrRequest.AuthToken)
	if err != nil {
		return err
	}
	defer file.DeleteDirectory(artifactsDir)

	note, err := d.p.Parse(mrRequest.ArtifactFormat, mrRequest.ArtifactFileName, artifactsDir, mrRequest.VulnerabilityMgmtId)
	if err != nil {
		return err
	}

	err = d.c.SendNote(note, mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.AuthToken)
	if err != nil {
		return err
	}

	log.Printf("%s Finished processing request: %v\n", time.Now().Format(time.DateTime), mrRequest)

	return nil
}

func (d *MRDecorator) DecorateCli(mrRequest *models.MRRequest) error {
	log.Printf("%s Started processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.JobId)
	artifactsDir := ""

	if mrRequest.FilePath == "" {
		artifactsDir, err := d.c.GetArtifact(mrRequest.ProjectId, mrRequest.JobId, mrRequest.ArtifactFileName, mrRequest.AuthToken)
		if err != nil {
			return err
		}
		defer file.DeleteDirectory(artifactsDir)
	} else {
		artifactsDir, mrRequest.ArtifactFileName = filepath.Split(mrRequest.FilePath)
	}

	note, err := d.p.Parse(mrRequest.ArtifactFormat, mrRequest.ArtifactFileName, artifactsDir, mrRequest.VulnerabilityMgmtId)
	if err != nil {
		return err
	}

	err = d.c.SendNote(note, mrRequest.ProjectId, mrRequest.MergeRequestIid, mrRequest.AuthToken)
	if err != nil {
		return err
	}

	log.Printf("%s Finished processing request: %v\n", time.Now().Format(time.DateTime), mrRequest)

	return nil
}
