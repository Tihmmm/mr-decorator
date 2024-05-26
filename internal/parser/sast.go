package parser

import (
	"github.com/Tihmmm/mr-decorator/internal/config"
)

type Sast interface {
	vulnCount() int
	ToGenSast(cfg config.SastParserConfig, vulnMgmtId int) GenSast
}

type GenSast struct {
	HighCount          int
	CriticalCount      int
	HcCount            int
	HighVulns          []Vulnerability
	CriticalVulns      []Vulnerability
	VulnMgmtProjectUrl string
	VulnMgmtReportPath string
}

type Vulnerability struct {
	Name             string
	Location         string
	VulnMgmtInstance string
}

func (sast *GenSast) ApplyLimit() {
	if sast.HighCount > sastVulnsDisplayed {
		sast.HighVulns = sast.HighVulns[:sastVulnsDisplayed]
	}
}
