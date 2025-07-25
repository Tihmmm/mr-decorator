package parser

import (
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/errors"
	"github.com/Tihmmm/mr-decorator/internal/models"
	"github.com/Tihmmm/mr-decorator/pkg/file"
	"github.com/Tihmmm/mr-decorator/pkg/templater"
	"log"
)

const (
	scaVulnsDisplayed  = 10
	sastVulnsDisplayed = scaVulnsDisplayed
)

type Parser interface {
	Parse(format string, fileName string, fileDir string, vulnMgmtId int) (note string, err error)
}

type ArtifactParser struct {
	cfg config.ParserConfig
}

func NewParser(cfg config.ParserConfig) Parser {
	parser := &ArtifactParser{
		cfg: cfg,
	}
	return parser
}

func (p *ArtifactParser) Parse(format string, fileName string, dir string, vulnMgmtId int) (string, error) {
	switch format {
	case models.FprFn:
		var f fpr
		if err := ParseFprFile(dir, &f); err != nil {
			return "", err
		}
		sast := f.ToGenSast(p.cfg.SastParserConfig, vulnMgmtId)
		sast.ApplyLimit()
		note, err := templater.ExecToString(baseTemplateSast, &sast)
		if err != nil {
			log.Printf("Error executing sast template: %s\n", err)
			return "", err
		}
		return note, nil
	case models.CyclonedxJsonFn:
		var dx cycloneDX
		if err := file.ParseJsonFile(dir, fileName, &dx); err != nil {
			return "", err
		}
		sca := dx.ToGenSca(p.cfg.ScaParserConfig, vulnMgmtId)
		sca.ApplyLimit()
		note, err := templater.ExecToString(baseTemplateSca, &sca)
		if err != nil {
			log.Printf("Error executing cdx template: %s\n", err)
			return "", err
		}
		return note, nil
	case models.DependencyCheckJsonFn:
		var dc dependencyCheck
		if err := file.ParseJsonFile(dir, fileName, &dc); err != nil {
			return "", err
		}
		sca := dc.ToGenSca(p.cfg.ScaParserConfig, vulnMgmtId)
		sca.ApplyLimit()
		note, err := templater.ExecToString(baseTemplateSca, &sca)
		if err != nil {
			log.Printf("Error executing depcheck template: %s\n", err)
			return "", err
		}
		return note, nil
	default:
		log.Printf("Invalid format: %s\n", format)
		return "", &errors.FormatError{}
	}
}
